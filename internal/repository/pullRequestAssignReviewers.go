package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/consts"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"
	"pr-reviewer-assignment-service/internal/model"
)

func (repo *repository) PullRequestAssignReviewers(ctx context.Context, params entity.PullRequestAssignReviewersParams) (entity.PullRequestAssignReviewersResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}
	defer tx.Rollback()

	prExists, err := doesPullRequestExist(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}
	if !prExists {
		return entity.PullRequestAssignReviewersResult{FoundPullRequest: false}, nil
	}

	prStatus, err := getPullRequestStatus(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}
	if prStatus == enum.PR_STATUS_MERGED {
		return entity.PullRequestAssignReviewersResult{
			FoundPullRequest:    true,
			IsPullRequestMerged: true,
		}, nil
	}

	numberOfCandidatesToAssign, err := getNumberOfCandidatesToAssign(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}
	if numberOfCandidatesToAssign == 0 {
		return entity.PullRequestAssignReviewersResult{
			FoundPullRequest:      true,
			IsPullRequestMerged:   false,
			HasMaxReviewersAmount: true,
		}, nil
	}

	assignedCnt, err := assignNewReviewers(tx, params.PullRequestId, numberOfCandidatesToAssign)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}
	if assignedCnt == 0 {
		return entity.PullRequestAssignReviewersResult{
			FoundPullRequest:      true,
			IsPullRequestMerged:   false,
			HasMaxReviewersAmount: false,
			FoundCandidate:        false,
		}, nil
	}

	resultPR, err := getIncompletePRById(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.PullRequestAssignReviewersResult{}, err
	}

	return entity.PullRequestAssignReviewersResult{
		FoundPullRequest:      true,
		IsPullRequestMerged:   false,
		HasMaxReviewersAmount: false,
		FoundCandidate:        true,
		PullRequestIncomplete: convertModelToEntity_OneIncompletePR(resultPR),
	}, nil
}

func getNumberOfCandidatesToAssign(tx *sql.Tx, prId string) (int, error) {
	currentNumberOfAssignedReviewers, err := getNumberOfAssignedReviewers(tx, prId)
	if err != nil {
		return 0, err
	}
	return consts.MAX_REVIEWERS_PER_PULL_REQUEST - currentNumberOfAssignedReviewers, nil
}

func getNumberOfAssignedReviewers(tx *sql.Tx, prId string) (int, error) {
	reviewerIds, err := getAssignedReviewersIds(tx, prId)
	if err != nil {
		return 0, err
	}
	return len(reviewerIds), nil
}

func assignNewReviewers(tx *sql.Tx, prId string, limit int) (assignedCnt int, err error) {
	candidateIds, err := getNewCandidateIds(tx, prId, limit)
	if err != nil {
		return 0, err
	}

	err = assignReviewersToPullRequest(tx, prId, candidateIds)
	if err != nil {
		return 0, err
	}
	return len(candidateIds), nil
}

func getNewCandidateIds(tx *sql.Tx, prId string, limit int) ([]string, error) {
	author, err := getPullRequestAuthor(tx, prId)
	if err != nil {
		return []string{}, err
	}

	candidates, err := findNewReviewCandidates(tx, author.TeamName, author.UserId, prId, limit)
	if err != nil {
		return []string{}, err
	}
	if len(candidates) == 0 {
		return []string{}, nil
	}

	candidateIds := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		candidateIds = append(candidateIds, candidate.UserId)
	}
	return candidateIds, nil
}

func getIncompletePRById(tx *sql.Tx, prId string) (model.PullRequestIncomplete, error) {
	prs, err := getIncompletePRsMainInfo(tx, []string{prId})
	if err != nil {
		return model.PullRequestIncomplete{}, err
	}

	pr := prs[0]
	currentPRReviewers, err := getPRReviewers(tx, prId)
	if err != nil {
		return model.PullRequestIncomplete{}, err
	}
	pr.AssignedReviewers = currentPRReviewers
	return pr, nil
}
