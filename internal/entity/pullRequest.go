package entity

type PullRequestCreateParams struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
}

type PullRequestCreateResult struct {
	PullRequest        PullRequest
	FoundAuthorAndTeam bool
	FoundPR            bool
}

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
	CreatedAt         string
	MergedAt          string
}
