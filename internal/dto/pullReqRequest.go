package dto

type PullRequestCreate_Request struct {
	PullRequestId   string `json:"pull_request_id" validate:"required"`
	PullRequestName string `json:"pull_request_name" validate:"required"`
	AuthorId        string `json:"author_id" validate:"required"`
}

type PullRequestMerge_Request struct {
	PullRequestId string `json:"pull_request_id" validate:"required"`
}

type PullRequestReassign_Request struct {
	PullRequestId string `json:"pull_request_id" validate:"required"`
	OldReviewerId string `json:"old_user_id" validate:"required"`
}
