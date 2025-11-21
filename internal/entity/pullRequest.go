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

type PullRequestMergeParams struct {
	PullRequestId string
}

type PullRequestMergeResult struct {
	PullRequest PullRequest
	FoundPR     bool
}

type PullRequestReassignParams struct {
	PullRequestId string
	OldReviewerId string
}

type PullRequestReassignResult struct {
	PullRequest           PullRequest
	NewReviewerId         string
	FoundPR               bool
	FoundOldReviewer      bool
	IsPullRequestMerged   bool
	IsOldReviewerAssigned bool
	FoundCandidate        bool
}

type PullRequestShort struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}
