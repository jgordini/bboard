package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)


func addNewComment(ctx context.Context, c *cmd.AddNewComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		isApproved := !tenant.IsModerationEnabled || !user.RequiresModeration()
		var id int
		if err := trx.Get(&id, `
			INSERT INTO comments (tenant_id, post_id, content, user_id, created_at, is_approved) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id
		`, tenant.ID, c.Post.ID, c.Content, user.ID, time.Now(), isApproved); err != nil {
			return errors.Wrap(err, "failed add new comment")
		}

		q := &query.GetCommentByID{CommentID: id}
		if err := getCommentByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result

		return nil
	})
}

func toggleCommentReaction(ctx context.Context, c *cmd.ToggleCommentReaction) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var added bool
		err := trx.Scalar(&added, `
			WITH toggle_reaction AS (
				INSERT INTO reactions (comment_id, user_id, emoji, created_on)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (comment_id, user_id, emoji) DO NOTHING
				RETURNING true AS added
			),
			delete_existing AS (
				DELETE FROM reactions
				WHERE comment_id = $1 AND user_id = $2 AND emoji = $3
				AND NOT EXISTS (SELECT 1 FROM toggle_reaction)
				RETURNING false AS added
			)
			SELECT COALESCE(
				(SELECT added FROM toggle_reaction),
				(SELECT added FROM delete_existing),
				false
			)
		`, c.Comment.ID, user.ID, c.Emoji, time.Now())

		if err != nil {
			return errors.Wrap(err, "failed to toggle reaction")
		}

		c.Result = added
		return nil
	})
}

func updateComment(ctx context.Context, c *cmd.UpdateComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			UPDATE comments SET content = $1, edited_at = $2, edited_by_id = $3 
			WHERE id = $4 AND tenant_id = $5`, c.Content, time.Now(), user.ID, c.CommentID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update comment")
		}
		return nil
	})
}

func deleteComment(ctx context.Context, c *cmd.DeleteComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if _, err := trx.Execute(
			"UPDATE comments SET deleted_at = $1, deleted_by_id = $2 WHERE id = $3 AND tenant_id = $4",
			time.Now(), user.ID, c.CommentID, tenant.ID,
		); err != nil {
			return errors.Wrap(err, "failed delete comment")
		}
		return nil
	})
}

func setCommentPinned(ctx context.Context, c *cmd.SetCommentPinned) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Pinned {
			_, err := trx.Execute(`
				UPDATE comments SET pinned_at = $1, pinned_by_id = $2 WHERE id = $3 AND tenant_id = $4`,
				time.Now(), user.ID, c.CommentID, tenant.ID,
			)
			if err != nil {
				return errors.Wrap(err, "failed to pin comment")
			}
		} else {
			_, err := trx.Execute(`
				UPDATE comments SET pinned_at = NULL, pinned_by_id = NULL WHERE id = $1 AND tenant_id = $2`,
				c.CommentID, tenant.ID,
			)
			if err != nil {
				return errors.Wrap(err, "failed to unpin comment")
			}
		}
		return nil
	})
}

func getCommentPostID(ctx context.Context, q *query.GetCommentPostID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = 0
		err := trx.Scalar(&q.Result, "SELECT post_id FROM comments WHERE id = $1 AND tenant_id = $2", q.CommentID, tenant.ID)
		if err != nil {
			return err
		}
		return nil
	})
}

func getCommentByID(ctx context.Context, q *query.GetCommentByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = nil

		comment := dbEntities.Comment{}
		err := trx.Get(&comment,
			`SELECT c.id, 
							c.content, 
							c.created_at, 
							c.edited_at, 
							c.is_approved,
							u.id AS user_id, 
							u.name AS user_name,
							u.email AS user_email,
							u.role AS user_role, 
							u.status AS user_status,
							u.avatar_type AS user_avatar_type,
							u.avatar_bkey AS user_avatar_bkey, 
							e.id AS edited_by_id, 
							e.name AS edited_by_name,
							e.email AS edited_by_email,
							e.role AS edited_by_role,
							e.status AS edited_by_status,
							e.avatar_type AS edited_by_avatar_type,
							e.avatar_bkey AS edited_by_avatar_bkey
			FROM comments c
			INNER JOIN users u
			ON u.id = c.user_id
			AND u.tenant_id = c.tenant_id
			LEFT JOIN users e
			ON e.id = c.edited_by_id
			AND e.tenant_id = c.tenant_id
			WHERE c.id = $1
			AND c.tenant_id = $2
			AND c.deleted_at IS NULL`, q.CommentID, tenant.ID)

		if err != nil {
			return err
		}

		q.Result = comment.ToModel(ctx)
		return nil
	})
}

func flagComment(ctx context.Context, c *cmd.FlagComment) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			INSERT INTO comment_flags (tenant_id, comment_id, user_id, created_at, reason)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (tenant_id, comment_id, user_id) DO NOTHING`,
			tenant.ID, c.CommentID, user.ID, time.Now(), c.Reason,
		)
		if err != nil {
			return errors.Wrap(err, "failed to flag comment")
		}
		return nil
	})
}

func getCommentFlagsCount(ctx context.Context, q *query.GetCommentFlagsCount) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		err := trx.Scalar(&q.Result, "SELECT COUNT(*) FROM comment_flags WHERE comment_id = $1 AND tenant_id = $2", q.CommentID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get comment flags count")
		}
		return nil
	})
}

func getCommentFlagsCountsForPost(ctx context.Context, q *query.GetCommentFlagsCountsForPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		type row struct {
			CommentID int `db:"comment_id"`
			Count     int `db:"count"`
		}
		var rows []*row
		err := trx.Select(&rows, `
			SELECT cf.comment_id, COUNT(*) AS count
			FROM comment_flags cf
			INNER JOIN comments c ON c.id = cf.comment_id AND c.tenant_id = cf.tenant_id
			WHERE c.post_id = $1 AND cf.tenant_id = $2
			GROUP BY cf.comment_id`,
			q.PostID, tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to get comment flags counts for post")
		}
		q.Result = make(map[int]int)
		for _, r := range rows {
			q.Result[r.CommentID] = r.Count
		}
		return nil
	})
}

func getFlaggedComments(ctx context.Context, q *query.GetFlaggedComments) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.FlaggedCommentItem, 0)
		if tenant == nil {
			return errors.New("tenant is required")
		}
		type row struct {
			CommentID  int       `db:"comment_id"`
			Content    string    `db:"content"`
			CreatedAt  time.Time `db:"created_at"`
			UserID     int       `db:"user_id"`
			UserName   string    `db:"user_name"`
			UserEmail  string    `db:"user_email"`
			UserRole   int       `db:"user_role"`
			UserStatus int       `db:"user_status"`
			AvatarType int       `db:"user_avatar_type"`
			AvatarBkey string    `db:"user_avatar_bkey"`
			PostNumber int       `db:"post_number"`
			PostTitle  string    `db:"post_title"`
			PostSlug   string    `db:"post_slug"`
			FlagsCount int       `db:"flags_count"`
		}
		var rows []*row
		err := trx.Select(&rows, `
			SELECT c.id AS comment_id, c.content, c.created_at,
				u.id AS user_id, u.name AS user_name, u.email AS user_email, u.role AS user_role, u.status AS user_status,
				u.avatar_type AS user_avatar_type, u.avatar_bkey AS user_avatar_bkey,
				p.number AS post_number, p.title AS post_title, p.slug AS post_slug,
				COUNT(cf.id)::int AS flags_count
			FROM comment_flags cf
			INNER JOIN comments c ON c.id = cf.comment_id AND c.tenant_id = cf.tenant_id AND c.deleted_at IS NULL
			INNER JOIN posts p ON p.id = c.post_id AND p.tenant_id = c.tenant_id
			INNER JOIN users u ON u.id = c.user_id AND u.tenant_id = c.tenant_id
			WHERE cf.tenant_id = $1
			GROUP BY c.id, c.content, c.created_at, u.id, u.name, u.email, u.role, u.status, u.avatar_type, u.avatar_bkey, p.number, p.title, p.slug
			ORDER BY flags_count DESC, c.created_at DESC`,
			tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to get flagged comments")
		}
		for _, r := range rows {
			userEnt := &entity.User{ID: r.UserID, Name: r.UserName, Email: r.UserEmail}
			userEnt.Role = enum.Role(r.UserRole)
			userEnt.Status = enum.UserStatus(r.UserStatus)
			item := &query.FlaggedCommentItem{
				Comment:    &entity.Comment{ID: r.CommentID, Content: r.Content, CreatedAt: r.CreatedAt, User: userEnt},
				PostNumber: r.PostNumber,
				PostTitle:  r.PostTitle,
				PostSlug:   r.PostSlug,
				FlagsCount: r.FlagsCount,
			}
			q.Result = append(q.Result, item)
		}
		return nil
	})
}

func clearCommentFlags(ctx context.Context, c *cmd.ClearCommentFlags) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`DELETE FROM comment_flags WHERE comment_id = $1 AND tenant_id = $2`,
			c.CommentID, tenant.ID,
		)
		if err != nil {
			return errors.Wrap(err, "failed to clear comment flags")
		}
		return nil
	})
}

func getCommentsByPost(ctx context.Context, q *query.GetCommentsByPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*entity.Comment, 0)
		if tenant == nil {
			return errors.New("tenant is required")
		}

		comments := []*dbEntities.Comment{}
		userId := 0
		if user != nil {
			userId = user.ID
		}
		
		// Build approval filter based on user permissions
		approvalFilter := ""
		if user != nil && user.IsCollaborator() {
			// Admins and collaborators can see all comments
			approvalFilter = ""
		} else if user != nil {
			// Regular users can see approved comments + their own unapproved comments
			approvalFilter = fmt.Sprintf(" AND (c.is_approved = true OR c.user_id = %d)", user.ID)
		} else {
			// Anonymous users can only see approved comments
			approvalFilter = " AND c.is_approved = true"
		}
		
		query := fmt.Sprintf(`
			WITH agg_attachments AS ( 
					SELECT 
							c.id as comment_id, 
							ARRAY_REMOVE(ARRAY_AGG(at.attachment_bkey), NULL) as attachment_bkeys
					FROM attachments at
					INNER JOIN comments c
					ON at.tenant_id = c.tenant_id
					AND at.post_id = c.post_id
					AND at.comment_id = c.id
					WHERE at.post_id = $1
					AND at.tenant_id = $2
					AND at.comment_id IS NOT NULL
					GROUP BY c.id 
			),
			agg_reactions AS (
				SELECT 
					comment_id,
					json_agg(json_build_object(
						'emoji', emoji,
						'count', count,
						'includesMe', CASE WHEN $3 = ANY(user_ids) THEN true ELSE false END
					) ORDER BY count DESC) as reaction_counts
				FROM (
					SELECT 
						comment_id, 
						emoji, 
						COUNT(*) as count,
						array_agg(user_id) as user_ids
					FROM reactions
					WHERE comment_id IN (SELECT id FROM comments WHERE post_id = $1)
					GROUP BY comment_id, emoji
				) r
				GROUP BY comment_id
			)
			SELECT c.id, 
					c.content, 
					c.created_at, 
					c.edited_at, 
					c.is_approved,
					c.pinned_at,
					u.id AS user_id, 
					u.name AS user_name,
					u.email AS user_email,
					u.role AS user_role, 
					u.status AS user_status, 
					u.avatar_type AS user_avatar_type, 
					u.avatar_bkey AS user_avatar_bkey, 
					e.id AS edited_by_id, 
					e.name AS edited_by_name,
					e.email AS edited_by_email,
					e.role AS edited_by_role,
					e.status AS edited_by_status,
					e.avatar_type AS edited_by_avatar_type, 
					e.avatar_bkey AS edited_by_avatar_bkey,
					pinner.id AS pinned_by_id,
					pinner.name AS pinned_by_name,
					pinner.email AS pinned_by_email,
					pinner.role AS pinned_by_role,
					pinner.status AS pinned_by_status,
					pinner.avatar_type AS pinned_by_avatar_type,
					pinner.avatar_bkey AS pinned_by_avatar_bkey,
					at.attachment_bkeys,
					ar.reaction_counts
			FROM comments c
			INNER JOIN posts p
			ON p.id = c.post_id
			AND p.tenant_id = c.tenant_id
			INNER JOIN users u
			ON u.id = c.user_id
			AND u.tenant_id = c.tenant_id
			LEFT JOIN users e
			ON e.id = c.edited_by_id
			AND e.tenant_id = c.tenant_id
			LEFT JOIN users pinner
			ON pinner.id = c.pinned_by_id
			AND pinner.tenant_id = c.tenant_id
			LEFT JOIN agg_attachments at
			ON at.comment_id = c.id
			LEFT JOIN agg_reactions ar
			ON ar.comment_id = c.id
			WHERE p.id = $1
			AND p.tenant_id = $2
			AND c.deleted_at IS NULL%s
			ORDER BY c.pinned_at DESC NULLS LAST, c.created_at DESC`, approvalFilter)
		
		err := trx.Select(&comments, query, q.Post.ID, tenant.ID, userId)
		if err != nil {
			return errors.Wrap(err, "failed get comments of post with id '%d'", q.Post.ID)
		}

		q.Result = make([]*entity.Comment, len(comments))
		for i, comment := range comments {
			q.Result[i] = comment.ToModel(ctx)
		}
		return nil
	})
}
