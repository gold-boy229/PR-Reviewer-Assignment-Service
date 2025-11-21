package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
)

func (repo *repository) GetUserAssignedPullRequests(ctx context.Context, params entity.UserGetAssignedPRParams) (entity.UserGetAssignedPRResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.UserGetAssignedPRResult{}, err
	}
	defer tx.Rollback()

	shortPRs, err := GetReviewerAssignedPullRequests(tx, params.UserId)
	if err != nil {
		return entity.UserGetAssignedPRResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.UserGetAssignedPRResult{}, err
	}

	return entity.UserGetAssignedPRResult{
		UserId:       params.UserId,
		PullRequests: convertModelToEntity_ManyShortPRs(shortPRs),
	}, nil
}

func GetReviewerAssignedPullRequests(tx *sql.Tx, reviewerId string) ([]model.PullRequestShort, error) {
	resultNumber, err := GetReviewerAssignedPRsNumber(tx, reviewerId)
	if err != nil {
		return []model.PullRequestShort{}, err
	}
	resultPRs := make([]model.PullRequestShort, 0, resultNumber)
	var currentPR model.PullRequestShort

	query := `	SELECT 
					pull_requests.pull_request_id,
					pull_requests.pull_request_name,
					pull_requests.author_id,
					pull_requests.status
				FROM pull_requests__M2M__users
				INNER JOIN pull_requests
					ON pull_requests__M2M__users.pull_request_id = pull_requests.pull_request_id
				WHERE pull_requests__M2M__users.reviewer_id = $1`
	rows, err := tx.Query(query, reviewerId)
	if err != nil {
		return []model.PullRequestShort{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&currentPR.PullRequestId,
			&currentPR.PullRequestName,
			&currentPR.AuthorId,
			&currentPR.Status,
		)
		if err != nil {
			return []model.PullRequestShort{}, err
		}
		resultPRs = append(resultPRs, currentPR)
	}

	err = rows.Err()
	if err != nil {
		return []model.PullRequestShort{}, err
	}
	return resultPRs, nil
}

func GetReviewerAssignedPRsNumber(tx *sql.Tx, reviewerId string) (int, error) {
	query := `	SELECT COUNT(*)
				FROM pull_requests__M2M__users
				WHERE reviewer_id = $1`
	var resultNumber int
	err := tx.QueryRow(query, reviewerId).Scan(&resultNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return resultNumber, nil
}

func convertModelToEntity_ManyShortPRs(shortPRs []model.PullRequestShort) []entity.PullRequestShort {
	result := make([]entity.PullRequestShort, 0, len(shortPRs))
	for _, pr := range shortPRs {
		result = append(result, convertModelToEntity_OneShortPR(pr))
	}
	return result
}

func convertModelToEntity_OneShortPR(pr model.PullRequestShort) entity.PullRequestShort {
	return entity.PullRequestShort{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
	}
}
