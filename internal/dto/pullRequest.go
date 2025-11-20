package dto

type PullRequestCreate_Request struct {
	PullRequestId   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorId        string `json:"author_id" validate:"required"`
}

type PullRequestCreate_Response struct {
	PullRequestId        string   `json:"pull_request_id"`
	PullRequestName      string   `json:"pull_request_name"`
	AuthorId             string   `json:"author_id"`
	Status               string   `json:"status"`
	AssignedReviewersIds []string `json:"assigned_reviewers"`
	CreatedAt            string   `json:"createdAt"`
	MergedAt             string   `json:"mergedAt,omitempty"`
}
