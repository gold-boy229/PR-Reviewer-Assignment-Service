package dto

type PullRequestCreate_Response struct {
	PullRequest_Response PullRequest_Response `json:"pr"`
}

type PullRequestMerge_Response struct {
	PullRequest_Response PullRequest_Response `json:"pr"`
}

type PullRequestReassign_Response struct {
	PullRequest_Response PullRequest_Response `json:"pr"`
	NewReviewerId        string               `json:"replaced_by"`
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

type PullRequestShort_Response struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          string `json:"status"`
}

type PullRequestGetIncomplete_Response struct {
	TeamName      string                           `json:"team_name"`
	IncompletePRs []PullRequestIncomplete_Response `json:"incomplete_prs"`
}
type PullRequestIncomplete_Response struct {
	PullRequestId     string                `json:"pull_request_id"`
	PullRequestName   string                `json:"pull_request_name"`
	AuthorId          string                `json:"author_id"`
	Status            string                `json:"status"`
	AssignedReviewers []TeamMember_Response `json:"assigned_reviewers"`
	CreatedAt         string                `json:"createdAt"`
}
