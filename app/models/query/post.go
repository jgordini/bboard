package query

import (
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
)

type PostIsReferenced struct {
	PostID int

	Result bool
}

type CountPostPerStatus struct {
	Result map[enum.PostStatus]int
}

type GetPostByID struct {
	PostID int

	Result *entity.Post
}

type GetPostBySlug struct {
	Slug string

	Result *entity.Post
}

type GetPostByNumber struct {
	Number int

	Result *entity.Post
}

type SearchPosts struct {
	Query            string
	View             string
	Limit            string
	Statuses         []enum.PostStatus
	Tags             []string
	MyVotesOnly      bool
	NoTagsOnly       bool
	MyPostsOnly      bool
	ModerationFilter string // "pending", "approved", or empty (all)

	Result []*entity.Post
}

type FindSimilarPosts struct {
	Query string

	Result []*entity.Post
}

type GetAllPosts struct {
	Result []*entity.Post
}

// LeaderboardPost is a post with vote count for leaderboard display
type LeaderboardPost struct {
	Number     int    `json:"number" db:"number"`
	Title      string `json:"title" db:"title"`
	Slug       string `json:"slug" db:"slug"`
	VotesCount int    `json:"votesCount" db:"votes_count"`
	UserID     int    `json:"userId" db:"user_id"`
	UserName   string `json:"userName" db:"user_name"`
}

// LeaderboardUser is a user with total votes received on their ideas
type LeaderboardUser struct {
	UserID     int    `json:"userId" db:"user_id"`
	UserName   string `json:"userName" db:"user_name"`
	VotesCount int    `json:"votesCount" db:"votes_count"`
}

// GetTopPostsByVotes returns posts ordered by votes count (for leaderboard)
type GetTopPostsByVotes struct {
	Limit  int
	Result []*LeaderboardPost
}

// GetTopUsersByVotes returns users ranked by total votes received on their posts (for leaderboard)
type GetTopUsersByVotes struct {
	Limit  int
	Result []*LeaderboardUser
}

func (q *SearchPosts) SetStatusesFromStrings(statuses []string) {
	for _, v := range statuses {
		var postStatus enum.PostStatus
		if err := postStatus.UnmarshalText([]byte(v)); err == nil {
			q.Statuses = append(q.Statuses, postStatus)
		}
	}
}

// GetPostFlagsCount returns the number of flags on a specific post
type GetPostFlagsCount struct {
	PostID int
	Result int
}

// FlaggedPostItem represents a flagged post with details for admin view
type FlaggedPostItem struct {
	Post       *entity.Post `json:"post"`
	FlagsCount int          `json:"flagsCount"`
}

// GetFlaggedPosts returns all posts that have been flagged (for admin moderation view)
type GetFlaggedPosts struct {
	Result []*FlaggedPostItem
}
