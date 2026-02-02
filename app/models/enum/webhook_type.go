package enum

// WebhookType is the type of a webhook
type WebhookType int

const (
	// WebhookNewPost is triggered on new post
	WebhookNewPost WebhookType = 1
	// WebhookNewComment is triggered on new comment
	WebhookNewComment WebhookType = 2
	// WebhookChangeStatus is triggered on post status change
	WebhookChangeStatus WebhookType = 3
	// WebhookDeletePost is triggered on post deletion
	WebhookDeletePost WebhookType = 4
	// WebhookCommentFlagged is triggered when a comment is flagged for inappropriateness
	WebhookCommentFlagged WebhookType = 5
	// WebhookCommentPinned is triggered when a moderator pins or unpins a comment
	WebhookCommentPinned WebhookType = 6
	// WebhookCommentDeleted is triggered when a comment is removed (by author or moderator)
	WebhookCommentDeleted WebhookType = 7
	// WebhookPostFlagged is triggered when a post is flagged for inappropriateness
	WebhookPostFlagged WebhookType = 8
)

var webhookTypeIDs = map[WebhookType]string{
	WebhookNewPost:        "new_post",
	WebhookNewComment:     "new_comment",
	WebhookChangeStatus:   "change_status",
	WebhookDeletePost:     "delete_post",
	WebhookCommentFlagged: "comment_flagged",
	WebhookCommentPinned:  "comment_pinned",
	WebhookCommentDeleted: "comment_deleted",
	WebhookPostFlagged:    "post_flagged",
}

var webhookTypeName = map[string]WebhookType{
	"new_post":         WebhookNewPost,
	"new_comment":      WebhookNewComment,
	"change_status":    WebhookChangeStatus,
	"delete_post":      WebhookDeletePost,
	"comment_flagged":  WebhookCommentFlagged,
	"comment_pinned":   WebhookCommentPinned,
	"comment_deleted":  WebhookCommentDeleted,
	"post_flagged":     WebhookPostFlagged,
}

// MarshalText returns the Text version of the webhook type
func (t WebhookType) MarshalText() ([]byte, error) {
	return []byte(webhookTypeIDs[t]), nil
}

// UnmarshalText parse string into a webhook type
func (t *WebhookType) UnmarshalText(text []byte) error {
	*t = webhookTypeName[string(text)]
	return nil
}

// Name returns the name of a webhook status
func (t WebhookType) Name() string {
	name, ok := webhookTypeIDs[t]
	if ok {
		return name
	}
	return "unknown"
}
