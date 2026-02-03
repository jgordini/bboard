package cmd

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type AddNewPost struct {
	Title       string
	Description string

	Result *entity.Post
}

type UpdatePost struct {
	Post        *entity.Post
	Title       string
	Description string

	Result *entity.Post
}

type SetPostResponse struct {
	Post   *entity.Post
	Text   string
	Status enum.PostStatus
}

// SetPostPinned sets or clears the pinned state of a post (staff only)
type SetPostPinned struct {
	PostID int
	Pinned bool
}

// FlagPost records that a user flagged a post (idempotent: one flag per user per post)
type FlagPost struct {
	PostID int
	Reason string
}

// ClearPostFlags removes all flags from a post (admin only)
type ClearPostFlags struct {
	PostID int
}
