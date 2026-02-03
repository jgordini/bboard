package cmd

import (
	"github.com/getfider/fider/app/models/entity"
)

type AddNewComment struct {
	Post    *entity.Post
	Content string

	Result *entity.Comment
}

type UpdateComment struct {
	CommentID int
	Content   string
}

type DeleteComment struct {
	CommentID int
}

// FlagComment records that a user flagged a comment (idempotent: one flag per user per comment)
type FlagComment struct {
	CommentID int
	Reason    string
}

// ClearCommentFlags removes all flags from a comment (admin only)
type ClearCommentFlags struct {
	CommentID int
}

// SetCommentPinned sets or clears the pinned state of a comment (moderators only)
type SetCommentPinned struct {
	CommentID int
	Pinned    bool
}
