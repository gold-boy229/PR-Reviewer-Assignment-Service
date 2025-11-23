package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/consts"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"
	"pr-reviewer-assignment-service/internal/model"
)

func (repo *repository) PullRequestReassign(ctx context.Context, params entity.PullRequestReassignParams) (entity.PullRequestReassignResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	defer tx.Rollback()

	pullRequestExists, err := doesPullRequestExist(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	if !pullRequestExists {
		return entity.PullRequestReassignResult{FoundPR: false}, nil
	}

	oldReviewerExists, err := doesUserExists(tx, params.OldReviewerId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	if !oldReviewerExists {
		return entity.PullRequestReassignResult{
			FoundPR:          true,
			FoundOldReviewer: false,
		}, nil
	}

	prStatus, err := getPullRequestStatus(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	if prStatus == enum.PR_STATUS_MERGED {
		return entity.PullRequestReassignResult{
			FoundPR:             true,
			FoundOldReviewer:    true,
			IsPullRequestMerged: true,
		}, nil
	}

	isOldReviewerAssigned, err := isReviewerAssigned(tx, params.PullRequestId, params.OldReviewerId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	if !isOldReviewerAssigned {
		return entity.PullRequestReassignResult{
			FoundPR:               true,
			FoundOldReviewer:      true,
			IsPullRequestMerged:   false,
			IsOldReviewerAssigned: false,
		}, nil
	}

	prAuthor, err := getPullRequestAuthor(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}

	candidate, foundCandidate, err := findNewReviewer(tx, prAuthor.TeamName, prAuthor.UserId, params.PullRequestId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}
	if !foundCandidate {
		return entity.PullRequestReassignResult{
			FoundPR:               true,
			FoundOldReviewer:      true,
			IsPullRequestMerged:   false,
			IsOldReviewerAssigned: true,
			FoundCandidate:        false,
		}, nil
	}

	err = reassignReviewer(tx, params.PullRequestId, params.OldReviewerId, candidate.UserId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}

	pullRequest, err := getPullRequestById(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.PullRequestReassignResult{}, err
	}

	return entity.PullRequestReassignResult{
		PullRequest:           convertModelToEntity_PullRequest(pullRequest),
		NewReviewerId:         candidate.UserId,
		FoundPR:               true,
		FoundOldReviewer:      true,
		IsPullRequestMerged:   false,
		IsOldReviewerAssigned: true,
		FoundCandidate:        true,
	}, nil
}

func getPullRequestStatus(tx *sql.Tx, prId string) (string, error) {
	query := `	SELECT status
				FROM pull_requests
				WHERE pull_request_id = $1`
	var status string
	err := tx.QueryRow(query, prId).Scan(&status)
	return status, err
}

func isReviewerAssigned(tx *sql.Tx, pullRequestId, reviewerId string) (bool, error) {
	query := `	SELECT pull_request_id, reviewer_id
				FROM pull_requests__M2M__users
				WHERE pull_request_id = $1
					AND reviewer_id = $2`
	var res_prId, res_reviewerId string
	err := tx.QueryRow(query, pullRequestId, reviewerId).Scan(&res_prId, &res_reviewerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func getPullRequestAuthor(tx *sql.Tx, prId string) (model.User, error) {
	query := `	SELECT
					users.user_id,
					users.username,
					users.team_name,
					users.is_active
				FROM pull_requests
				INNER JOIN users
					ON pull_requests.author_id = users.user_id
				WHERE pull_requests.pull_request_id = $1`
	var author model.User
	err := tx.QueryRow(query, prId).Scan(
		&author.UserId,
		&author.Username,
		&author.TeamName,
		&author.IsActive,
	)
	return author, err
}

func findNewReviewer(tx *sql.Tx, teamName, authorId, prId string) (candidate model.User, found bool, err error) {
	const FIND_ONE_CANDIDATE int = 1
	candidates, err := findNewReviewCandidates(tx, teamName, authorId, prId, FIND_ONE_CANDIDATE)
	if err != nil {
		return model.User{}, false, err
	}
	if len(candidates) == 0 {
		return model.User{}, false, nil
	}
	return candidates[0], true, nil
}

func findNewReviewCandidates(tx *sql.Tx, teamName, authorId, prId string, limit int) ([]model.User, error) {
	query := `	SELECT 
					user_id,
					username,
					team_name,
					is_active
				FROM users
				WHERE team_name = $1
					AND is_active = true
					AND user_id != $2
					AND user_id NOT IN (
						SELECT reviewer_id
						FROM pull_requests__M2M__users
						WHERE pull_request_id = $3
					)
				LIMIT $4`
	rows, err := tx.Query(query, teamName, authorId, prId, limit)
	if err != nil {
		return []model.User{}, err
	}
	defer rows.Close()

	resultCandidates := make([]model.User, 0, consts.MAX_REVIEWERS_PER_PULL_REQUEST)
	var currentCandidate model.User

	for rows.Next() {
		err = rows.Scan(
			&currentCandidate.UserId,
			&currentCandidate.Username,
			&currentCandidate.TeamName,
			&currentCandidate.IsActive,
		)
		if err != nil {
			return []model.User{}, err
		}
		resultCandidates = append(resultCandidates, currentCandidate)
	}

	err = rows.Err()
	if err != nil {
		return []model.User{}, err
	}
	return resultCandidates, nil
}

func reassignReviewer(tx *sql.Tx, prId, oldReviewerId, newReviewerId string) error {
	query := `	UPDATE pull_requests__M2M__users
				SET reviewer_id = $1
				WHERE pull_request_id = $2
					AND reviewer_id = $3`
	_, err := tx.Exec(query, newReviewerId, prId, oldReviewerId)
	return err
}
