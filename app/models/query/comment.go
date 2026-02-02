package query

import (
	"github.com/getfider/fider/app/models/entity"
)

type GetCommentByID struct {
	CommentID int

	Result *entity.Comment
}

type GetCommentsByPost struct {
	Post *entity.Post

	Result []*entity.Comment
}

// GetCommentFlagsCount returns the number of flags on a comment (for moderator view only)
type GetCommentFlagsCount struct {
	CommentID int
	Result    int
}

// GetCommentFlagsCountsForPost returns flag counts per comment for a post (for moderator view). PostID is posts.id.
type GetCommentFlagsCountsForPost struct {
	PostID int       // posts.id
	Result map[int]int // commentID -> count
}

// FlaggedCommentItem is a comment with flag count for the moderator flagged queue
type FlaggedCommentItem struct {
	Comment    *entity.Comment
	PostNumber int
	PostTitle  string
	PostSlug   string
	FlagsCount int
}

// GetFlaggedComments returns comments that have at least one flag (for moderator queue)
type GetFlaggedComments struct {
	Result []*FlaggedCommentItem
}

// GetCommentPostID returns the post_id (posts.id) for a comment; 0 if not found
type GetCommentPostID struct {
	CommentID int
	Result    int
}
