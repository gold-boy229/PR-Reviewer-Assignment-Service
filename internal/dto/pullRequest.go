package dto

type PullRequestCreate_Request struct {
	PullRequestId   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorId        string `json:"author_id" validate:"required"`
}

type PullRequestMegre_Request struct {
	PullRequestId string `json:"pull_request_id" validate:"required"`
}

type PullRequestReassign_Request struct {
	PullRequestId string `json:"pull_request_id" validate:"required"`
	OldReviewerId string `json:"old_user_id" validate:"required"`
}

type PullRequest_Response struct {
	PullRequestId        string   `json:"pull_request_id"`
	PullRequestName      string   `json:"pull_request_name"`
	AuthorId             string   `json:"author_id"`
	Status               string   `json:"status"`
	AssignedReviewersIds []string `json:"assigned_reviewers"`
	CreatedAt            string   `json:"createdAt"`
	MergedAt             string   `json:"mergedAt,omitempty"`
}

type PullRequestReassign_Response struct {
	PullRequest_Response PullRequest_Response `json:"pr"`
	NewReviewerId        string               `json:"replaced_by"`
}
