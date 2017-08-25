package actuator

// PullRequestEvent mirrors the JSON structure of the Github Pull-Request event which
// is described here: https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEvent struct {
	Action      string      `json:"action" binding:"required"`
	Number      int32       `json:"number" binding:"required"`
	PullRequest PullRequest `json:"pull_request"`
}

// PullRequest mirrors the JSON structure of a pull request.
type PullRequest struct {
	ID          int32           `json:"id"`
	URL         string          `json:"url"`
	Title       string          `json:"title"`
	Body        string          `json:"body"`
	State       string          `json:"state"`
	CommentsURL string          `json:"comments_url"`
	Head        PullRequestHead `json:"head"`
}

// PullRequestHead contains information about the Pull Request branch
type PullRequestHead struct {
	Ref   string `json:"ref"`
	SHA   string `json:"sha"`
	Label string `json:"label"`
}
