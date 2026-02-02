package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pemistahl/lingua-go"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/gosimple/slug"
	"github.com/lib/pq"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

var (
	sqlSelectPostsWhere = `	WITH
													agg_tags AS (
														SELECT
																post_id,
																ARRAY_REMOVE(ARRAY_AGG(tags.slug), NULL) as tags
														FROM post_tags
														INNER JOIN tags
														ON tags.ID = post_tags.TAG_ID
														AND tags.tenant_id = post_tags.tenant_id
														WHERE post_tags.tenant_id = $1
														%s
														GROUP BY post_id
													),
													agg_comments AS (
															SELECT
																	post_id,
																	COUNT(CASE WHEN comments.created_at > CURRENT_DATE - INTERVAL '30 days' AND comments.is_approved = true THEN 1 END) as recent,
																	COUNT(CASE WHEN comments.is_approved = true THEN 1 END) as all
															FROM comments
															INNER JOIN posts
															ON posts.id = comments.post_id
															AND posts.tenant_id = comments.tenant_id
															WHERE posts.tenant_id = $1
															AND comments.deleted_at IS NULL
															GROUP BY post_id
													),
													agg_votes AS (
															SELECT
															post_id,
																	COUNT(CASE WHEN post_votes.created_at > CURRENT_DATE - INTERVAL '30 days'  THEN 1 END) as recent,
																	COUNT(*) as all
															FROM post_votes
															INNER JOIN posts
															ON posts.id = post_votes.post_id
															AND posts.tenant_id = post_votes.tenant_id
															WHERE posts.tenant_id = $1
															GROUP BY post_id
													)
													SELECT p.id,
																p.number,
																p.title,
																p.slug,
																p.description,
																p.created_at,
																p.search,
																COALESCE(agg_s.all, 0) as votes_count,
																COALESCE(agg_c.all, 0) as comments_count,
																COALESCE(agg_s.recent, 0) AS recent_votes_count,
																COALESCE(agg_c.recent, 0) AS recent_comments_count,
																p.status,
																u.id AS user_id,
																u.name AS user_name,
																u.email AS user_email,
																u.role AS user_role,
																u.status AS user_status,
																u.avatar_type AS user_avatar_type,
																u.avatar_bkey AS user_avatar_bkey,
																p.response,
																p.response_date,
																r.id AS response_user_id,
																r.name AS response_user_name,
																r.email AS response_user_email,
																r.role AS response_user_role,
																r.status AS response_user_status,
																r.avatar_type AS response_user_avatar_type,
																r.avatar_bkey AS response_user_avatar_bkey,
																d.number AS original_number,
																d.title AS original_title,
																d.slug AS original_slug,
																d.status AS original_status,
																COALESCE(agg_t.tags, ARRAY[]::text[]) AS tags,
																COALESCE(%s, false) AS has_voted,
																p.is_approved,
																p.pinned_at,
																pinner.id AS pinned_by_id,
																pinner.name AS pinned_by_name,
																pinner.email AS pinned_by_email,
																pinner.role AS pinned_by_role,
																pinner.status AS pinned_by_status,
																pinner.avatar_type AS pinned_by_avatar_type,
																pinner.avatar_bkey AS pinned_by_avatar_bkey
													FROM posts p
													INNER JOIN users u
													ON u.id = p.user_id
													AND u.tenant_id = $1
													LEFT JOIN users r
													ON r.id = p.response_user_id
													AND r.tenant_id = $1
													LEFT JOIN users pinner
													ON pinner.id = p.pinned_by_id
													AND pinner.tenant_id = $1
													LEFT JOIN posts d
													ON d.id = p.original_id
													AND d.tenant_id = $1
													LEFT JOIN agg_comments agg_c
													ON agg_c.post_id = p.id
													LEFT JOIN agg_votes agg_s
													ON agg_s.post_id = p.id
													LEFT JOIN agg_tags agg_t
													ON agg_t.post_id = p.id
													WHERE p.status != ` + strconv.Itoa(int(enum.PostDeleted)) + ` AND %s`
)

func postIsReferenced(ctx context.Context, q *query.PostIsReferenced) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = false

		exists, err := trx.Exists(`
			SELECT 1 FROM posts p
			INNER JOIN posts o
			ON o.tenant_id = p.tenant_id
			AND o.id = p.original_id
			WHERE p.tenant_id = $1
			AND o.id = $2`, tenant.ID, q.PostID)
		if err != nil {
			return errors.Wrap(err, "failed to check if post is referenced")
		}

		q.Result = exists
		return nil
	})
}

func getTopPostsByVotes(ctx context.Context, q *query.GetTopPostsByVotes) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.LeaderboardPost, 0)
		limit := q.Limit
		if limit <= 0 {
			limit = 10
		}
		err := trx.Select(&q.Result, `
			SELECT p.number, p.title, p.slug, u.id AS user_id, u.name AS user_name,
				COALESCE(vc.cnt, 0)::int AS votes_count
			FROM posts p
			INNER JOIN users u ON u.id = p.user_id AND u.tenant_id = p.tenant_id
			LEFT JOIN (SELECT post_id, COUNT(*) AS cnt FROM post_votes GROUP BY post_id) vc ON vc.post_id = p.id
			WHERE p.tenant_id = $1 AND p.status != $2
			ORDER BY votes_count DESC, p.number DESC
			LIMIT $3`, tenant.ID, enum.PostDeleted, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get top posts by votes")
		}
		return nil
	})
}

func getTopUsersByVotes(ctx context.Context, q *query.GetTopUsersByVotes) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.LeaderboardUser, 0)
		limit := q.Limit
		if limit <= 0 {
			limit = 10
		}
		type row struct {
			UserID     int    `db:"user_id"`
			UserName   string `db:"user_name"`
			VotesCount int    `db:"votes_count"`
		}
		var rows []*row
		err := trx.Select(&rows, `
			SELECT p.user_id, u.name AS user_name, COALESCE(SUM(vc.cnt), 0)::int AS votes_count
			FROM posts p
			INNER JOIN users u ON u.id = p.user_id AND u.tenant_id = p.tenant_id
			LEFT JOIN (SELECT post_id, COUNT(*) AS cnt FROM post_votes GROUP BY post_id) vc ON vc.post_id = p.id
			WHERE p.tenant_id = $1 AND p.status != $2
			GROUP BY p.user_id, u.name
			ORDER BY votes_count DESC
			LIMIT $3`, tenant.ID, enum.PostDeleted, limit)
		if err != nil {
			return errors.Wrap(err, "failed to get top users by votes")
		}
		for _, r := range rows {
			q.Result = append(q.Result, &query.LeaderboardUser{UserID: r.UserID, UserName: r.UserName, VotesCount: r.VotesCount})
		}
		return nil
	})
}

func setPostResponse(ctx context.Context, c *cmd.SetPostResponse) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Status == enum.PostDuplicate {
			return errors.New("Use MarkAsDuplicate to change an post status to Duplicate")
		}

		respondedAt := time.Now()
		if c.Post.Status == c.Status && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		_, err := trx.Execute(`
		UPDATE posts
		SET response = $3, original_id = NULL, response_date = $4, response_user_id = $5, status = $6
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Text, respondedAt, user.ID, c.Status)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = c.Status
		c.Post.Response = &entity.PostResponse{
			Text:        c.Text,
			RespondedAt: respondedAt,
			User:        user,
		}
		return nil
	})
}

func markPostAsDuplicate(ctx context.Context, c *cmd.MarkPostAsDuplicate) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		respondedAt := time.Now()
		if c.Post.Status == enum.PostDuplicate && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		var users []*dbEntities.User
		err := trx.Select(&users, "SELECT user_id AS id FROM post_votes WHERE post_id = $1 AND tenant_id = $2", c.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get votes of post with id '%d'", c.Post.ID)
		}

		for _, u := range users {
			err := bus.Dispatch(ctx, &cmd.AddVote{Post: c.Original, User: u.ToModel(ctx)})
			if err != nil {
				return err
			}
		}

		_, err = trx.Execute(`
		UPDATE posts
		SET response = '', original_id = $3, response_date = $4, response_user_id = $5, status = $6
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Original.ID, respondedAt, user.ID, enum.PostDuplicate)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = enum.PostDuplicate
		c.Post.Response = &entity.PostResponse{
			RespondedAt: respondedAt,
			User:        user,
			Original: &entity.OriginalPost{
				Number: c.Original.Number,
				Title:  c.Original.Title,
				Slug:   c.Original.Slug,
				Status: c.Original.Status,
			},
		}
		return nil
	})
}

func countPostPerStatus(ctx context.Context, q *query.CountPostPerStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {

		type dbStatusCount struct {
			Status enum.PostStatus `db:"status"`
			Count  int             `db:"count"`
		}

		q.Result = make(map[enum.PostStatus]int)
		stats := []*dbStatusCount{}
		err := trx.Select(&stats, "SELECT status, COUNT(*) AS count FROM posts WHERE tenant_id = $1 GROUP BY status", tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to count posts per status")
		}

		for _, v := range stats {
			q.Result[v.Status] = v.Count
		}
		return nil
	})
}

func addNewPost(ctx context.Context, c *cmd.AddNewPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		isApproved := !tenant.IsModerationEnabled || !user.RequiresModeration()
		var id int
		// Detect language using lingua-go
		lang := detectPostLanguage(c.Title, c.Description)

		err := trx.Get(&id,
			`INSERT INTO posts (title, slug, number, description, tenant_id, user_id, created_at, status, is_approved, language)
			 VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM posts p WHERE p.tenant_id = $4), $3, $4, $5, $6, 0, $7, $8)
			 RETURNING id`, c.Title, slug.Make(c.Title), c.Description, tenant.ID, user.ID, time.Now(), isApproved, lang)
		if err != nil {
			return errors.Wrap(err, "failed add new post")
		}

		q := &query.GetPostByID{PostID: id}
		if err := getPostByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result

		if err := internalAddSubscriber(trx, q.Result, tenant, user, false); err != nil {
			return err
		}

		return nil
	})
}

func updatePost(ctx context.Context, c *cmd.UpdatePost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		// Detect language using lingua-go
		lang := detectPostLanguage(c.Title, c.Description)
		_, err := trx.Execute(`UPDATE posts SET title = $1, slug = $2, description = $3, language = $4
								 WHERE id = $5 AND tenant_id = $6`, c.Title, slug.Make(c.Title), c.Description, lang, c.Post.ID, tenant.ID)

		if err != nil {
			return errors.Wrap(err, "failed update post")
		}

		q := &query.GetPostByID{PostID: c.Post.ID}
		if err := getPostByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result
		return nil
	})
}

// detectPostLanguage uses lingua-go to detect the language of a post and maps it to a PostgreSQL tsvector config or 'simple'.
// All language mappings are centralized in app/models/enum/locale.go
func detectPostLanguage(title, description string) string {
	// Get all supported lingua languages from the centralized locale enum
	linguaLangs := enum.GetLinguaLanguages()
	detector := lingua.NewLanguageDetectorBuilder().FromLanguages(linguaLangs...).Build()
	text := strings.TrimSpace(title + " " + description)
	lang, exists := detector.DetectLanguageOf(text)
	if !exists {
		return "simple"
	}

	locale, ok := enum.GetLocaleByLinguaLanguage(lang)
	if !ok {
		return "simple"
	}

	// Map to PostgreSQL tsvector config using the locale's PostgresConfig
	return locale.PostgresConfig
}

func getPostByID(ctx context.Context, q *query.GetPostByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildSinglePostQuery(user, "p.tenant_id = $1 AND p.id = $2"), tenant.ID, q.PostID)
		if err != nil {
			return errors.Wrap(err, "failed to get post with id '%d'", q.PostID)
		}
		q.Result = post
		return nil
	})
}

func getPostBySlug(ctx context.Context, q *query.GetPostBySlug) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildSinglePostQuery(user, "p.tenant_id = $1 AND p.slug = $2"), tenant.ID, q.Slug)
		if err != nil {
			return errors.Wrap(err, "failed to get post with slug '%s'", q.Slug)
		}
		q.Result = post
		return nil
	})
}

func getPostByNumber(ctx context.Context, q *query.GetPostByNumber) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildSinglePostQuery(user, "p.tenant_id = $1 AND p.number = $2"), tenant.ID, q.Number)
		if err != nil {
			return errors.Wrap(err, "failed to get post with number '%d'", q.Number)
		}
		q.Result = post
		return nil
	})
}

func preprocessSearchQuery(query string) string {
	// Common noise words that don't add search value

	noiseWords := env.SearchNoiseWords()

	words := strings.Fields(strings.ToLower(query))
	var filteredWords []string

	for _, word := range words {
		isNoise := false
		for _, noise := range noiseWords {
			if word == noise {
				isNoise = true
				break
			}
		}
		if !isNoise && len(word) > 2 { // Also filter very short words
			filteredWords = append(filteredWords, word)
		}
	}

	return strings.Join(filteredWords, " ")
}

func findSimilarPosts(ctx context.Context, q *query.FindSimilarPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		innerQuery := buildPostQuery(user, "p.tenant_id = $1 AND p.status = ANY($2)", "")

		filteredQuery := preprocessSearchQuery(q.Query)

		var (
			posts []*dbEntities.Post
			err   error
		)

		if filteredQuery == "" {
			q.Result = make([]*entity.Post, 0)
		} else {
			tsConfig := MapLocaleToTSConfig(tenant.Locale)

			// Build tsquery with AND operator between words and prefix matching on each word
			// The search column already contains both language-specific and simple tsvectors
			tsQueryExpr := fmt.Sprintf("to_tsquery('%s', regexp_replace(regexp_replace($3, '\\\\s+', ':* & ', 'g'), '$', ':*'))", tsConfig)
			tsQuerySimple := "to_tsquery('simple', regexp_replace(regexp_replace($3, '\\\\s+', ':* & ', 'g'), '$', ':*'))"

			// Use ts_rank_cd (cover density ranking) for better relevance scoring
			// Query against the generated search column which combines language-specific and simple tsvectors
			score := fmt.Sprintf("ts_rank_cd(q.search, %s) + ts_rank_cd(q.search, %s)", tsQueryExpr, tsQuerySimple)

			// Match against the pre-computed search column
			whereParts := fmt.Sprintf(`q.search @@ %s OR q.search @@ %s`, tsQueryExpr, tsQuerySimple)

			sql := fmt.Sprintf(`
				SELECT * FROM (%s) AS q
				WHERE %s
				ORDER BY %s DESC
				LIMIT 5
			`, innerQuery, whereParts, score)
			err = trx.Select(&posts, sql, tenant.ID, pq.Array([]enum.PostStatus{
				enum.PostOpen,
				enum.PostStarted,
				enum.PostPlanned,
				enum.PostCompleted,
				enum.PostDeclined,
			}), ToTSQuery(SanitizeString(q.Query)))
		}
		if err != nil {
			return errors.Wrap(err, "failed to find similar posts")
		}

		q.Result = make([]*entity.Post, len(posts))
		for i, post := range posts {
			q.Result[i] = post.ToModel(ctx)
		}
		return nil
	})
}

func searchPosts(ctx context.Context, q *query.SearchPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		innerQuery := buildPostQuery(user, "p.tenant_id = $1 AND p.status = ANY($2)", q.ModerationFilter)

		if q.Tags == nil {
			q.Tags = []string{}
		}

		if q.Statuses == nil {
			q.Statuses = []enum.PostStatus{}
		}

		if q.Limit != "all" {
			if _, err := strconv.Atoi(q.Limit); err != nil {
				q.Limit = "30"
			}
		}

		var (
			posts []*dbEntities.Post
			err   error
		)
		if q.Query != "" {
			tsQuery := ToTSQuery(SanitizeString(q.Query))
			if tsQuery == "" {
				q.Result = make([]*entity.Post, 0)
				return nil
			}

			tsConfig := MapLocaleToTSConfig(tenant.Locale)

			tsQueryExpr := fmt.Sprintf("to_tsquery('%s', regexp_replace(regexp_replace($3, '\\\\s+', ':* & ', 'g'), '$', ':*'))", tsConfig)
			tsQuerySimple := "to_tsquery('simple', regexp_replace(regexp_replace($3, '\\\\s+', ':* & ', 'g'), '$', ':*'))"

			score := fmt.Sprintf("ts_rank_cd(q.search, %s) + ts_rank_cd(q.search, %s)", tsQueryExpr, tsQuerySimple)

			whereParts := fmt.Sprintf(`q.search @@ %s OR q.search @@ %s`, tsQueryExpr, tsQuerySimple)

			sql := fmt.Sprintf(`
				SELECT * FROM (%s) AS q
				WHERE %s
				ORDER BY q.pinned_at DESC NULLS LAST, %s DESC
				LIMIT %s
			`, innerQuery, whereParts, score, q.Limit)
			err = trx.Select(&posts, sql, tenant.ID, pq.Array([]enum.PostStatus{
				enum.PostOpen,
				enum.PostStarted,
				enum.PostPlanned,
				enum.PostCompleted,
				enum.PostDeclined,
			}), tsQuery)
		} else {
			condition, statuses, sort := getViewData(*q)

			if q.MyPostsOnly {
				condition += " AND user_id = " + strconv.Itoa(user.ID)
			}

			// Second sort key: simple column (e.g. "id") -> "q.id"; expression (e.g. trending formula) -> qualify columns with "q."
			sortKey := sort
			if strings.Contains(sort, "(") {
				sortKey = strings.ReplaceAll(sort, "recent_votes_count", "q.recent_votes_count")
				sortKey = strings.ReplaceAll(sortKey, "recent_comments_count", "q.recent_comments_count")
				sortKey = strings.ReplaceAll(sortKey, "created_at", "q.created_at")
			} else {
				sortKey = "q." + sort
			}

			sql := fmt.Sprintf(`
				SELECT * FROM (%s) AS q
				WHERE 1 = 1 %s
				ORDER BY q.pinned_at DESC NULLS LAST, %s DESC
				LIMIT %s
			`, innerQuery, condition, sortKey, q.Limit)
			params := []interface{}{tenant.ID, pq.Array(statuses)}
			if len(q.Tags) > 0 {
				params = append(params, pq.Array(q.Tags))
			}
			err = trx.Select(&posts, sql, params...)
		}

		if err != nil {
			return errors.Wrap(err, "failed to search posts")
		}

		q.Result = make([]*entity.Post, len(posts))
		for i, post := range posts {
			q.Result[i] = post.ToModel(ctx)
		}
		return nil
	})
}

func getAllPosts(ctx context.Context, q *query.GetAllPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		searchQuery := &query.SearchPosts{View: "all", Limit: "all"}
		if err := searchPosts(ctx, searchQuery); err != nil {
			return errors.Wrap(err, "failed to get all posts")
		}
		q.Result = searchQuery.Result
		return nil
	})
}

func querySinglePost(ctx context.Context, trx *dbx.Trx, query string, args ...any) (*entity.Post, error) {
	post := dbEntities.Post{}

	if err := trx.Get(&post, query, args...); err != nil {
		return nil, err
	}

	return post.ToModel(ctx), nil
}

func buildPostQuery(user *entity.User, filter string, moderationFilter string) string {
	tagCondition := `AND tags.is_public = true`
	if user != nil && user.IsCollaborator() {
		tagCondition = ``
	}
	hasVotedSubQuery := "null"
	if user != nil {
		hasVotedSubQuery = fmt.Sprintf("(SELECT true FROM post_votes WHERE post_id = p.id AND user_id = %d)", user.ID)
	}

	// Add approval filtering based on moderation filter and user permissions
	approvalFilter := ""

	// If user is a collaborator and has specified a moderation filter, apply it
	if user != nil && user.IsCollaborator() && moderationFilter != "" {
		switch moderationFilter {
		case "pending":
			// Show only unapproved posts
			approvalFilter = " AND p.is_approved = false"
		case "approved":
			// Show only approved posts
			approvalFilter = " AND p.is_approved = true"
		}
		// If moderationFilter is neither "pending" nor "approved", show all posts (no filter)
	} else if user != nil {
		// Regular authenticated users can see approved posts + their own unapproved posts
		approvalFilter = fmt.Sprintf(" AND (p.is_approved = true OR p.user_id = %d)", user.ID)
	} else {
		// Anonymous users can only see approved posts
		approvalFilter = " AND p.is_approved = true"
	}

	combinedFilter := filter + approvalFilter
	return fmt.Sprintf(sqlSelectPostsWhere, tagCondition, hasVotedSubQuery, combinedFilter)
}

// buildSinglePostQuery is used for fetching individual posts (by ID, slug, or number)
// Collaborators can view any post for moderation purposes
func buildSinglePostQuery(user *entity.User, filter string) string {
	tagCondition := `AND tags.is_public = true`
	if user != nil && user.IsCollaborator() {
		tagCondition = ``
	}
	hasVotedSubQuery := "null"
	if user != nil {
		hasVotedSubQuery = fmt.Sprintf("(SELECT true FROM post_votes WHERE post_id = p.id AND user_id = %d)", user.ID)
	}

	// Approval filtering for single post views
	approvalFilter := ""
	if user != nil && user.IsCollaborator() {
		// Collaborators can view any post (for moderation purposes)
		approvalFilter = ""
	} else if user != nil {
		// Regular authenticated users can see approved posts + their own unapproved posts
		approvalFilter = fmt.Sprintf(" AND (p.is_approved = true OR p.user_id = %d)", user.ID)
	} else {
		// Anonymous users can only see approved posts
		approvalFilter = " AND p.is_approved = true"
	}

	combinedFilter := filter + approvalFilter
	return fmt.Sprintf(sqlSelectPostsWhere, tagCondition, hasVotedSubQuery, combinedFilter)
}

func setPostPinned(ctx context.Context, c *cmd.SetPostPinned) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Pinned {
			_, err := trx.Execute(`
				UPDATE posts SET pinned_at = $3, pinned_by_id = $4 WHERE id = $1 AND tenant_id = $2`,
				c.PostID, tenant.ID, time.Now(), user.ID)
			if err != nil {
				return errors.Wrap(err, "failed to pin post")
			}
		} else {
			_, err := trx.Execute(`
				UPDATE posts SET pinned_at = NULL, pinned_by_id = NULL WHERE id = $1 AND tenant_id = $2`,
				c.PostID, tenant.ID)
			if err != nil {
				return errors.Wrap(err, "failed to unpin post")
			}
		}
		return nil
	})
}

func flagPost(ctx context.Context, c *cmd.FlagPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`
			INSERT INTO post_flags (tenant_id, post_id, user_id, created_at, reason)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (tenant_id, post_id, user_id) DO NOTHING`,
			tenant.ID, c.PostID, user.ID, time.Now(), c.Reason,
		)
		if err != nil {
			return errors.Wrap(err, "failed to flag post")
		}
		return nil
	})
}

func getPostFlagsCount(ctx context.Context, q *query.GetPostFlagsCount) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		err := trx.Scalar(&q.Result, "SELECT COUNT(*) FROM post_flags WHERE post_id = $1 AND tenant_id = $2", q.PostID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get post flags count")
		}
		return nil
	})
}

func getFlaggedPosts(ctx context.Context, q *query.GetFlaggedPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = make([]*query.FlaggedPostItem, 0)
		if tenant == nil {
			return errors.New("tenant is required")
		}

		// First get posts with their flag counts
		type row struct {
			PostID     int `db:"post_id"`
			FlagsCount int `db:"flags_count"`
		}
		var rows []*row
		err := trx.Select(&rows, `
			SELECT p.id AS post_id, COUNT(pf.id)::int AS flags_count
			FROM post_flags pf
			INNER JOIN posts p ON p.id = pf.post_id AND p.tenant_id = pf.tenant_id AND p.status != $2
			WHERE pf.tenant_id = $1
			GROUP BY p.id
			ORDER BY flags_count DESC, p.created_at DESC`,
			tenant.ID, enum.PostDeleted,
		)
		if err != nil {
			return errors.Wrap(err, "failed to get flagged posts")
		}

		// Now fetch full post details for each flagged post
		for _, r := range rows {
			getPost := &query.GetPostByID{PostID: r.PostID}
			if err := bus.Dispatch(ctx, getPost); err != nil {
				continue // Skip posts that can't be loaded
			}
			q.Result = append(q.Result, &query.FlaggedPostItem{
				Post:       getPost.Result,
				FlagsCount: r.FlagsCount,
			})
		}
		return nil
	})
}
