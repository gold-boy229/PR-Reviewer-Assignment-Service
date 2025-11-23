package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/consts"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"
	"pr-reviewer-assignment-service/internal/model"

	"github.com/lib/pq"
)

func (repo *repository) PullRequestGetOpenIncompletePRs(ctx context.Context, params entity.PullRequestGetIncompleteParams) (entity.PullRequestGetIncompleteResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.PullRequestGetIncompleteResult{}, err
	}
	defer tx.Rollback()

	teamExists, err := doesTeamExist(tx, params.TeamName)
	if err != nil {
		return entity.PullRequestGetIncompleteResult{}, err
	}
	if !teamExists {
		return entity.PullRequestGetIncompleteResult{FoundTeam: false}, err
	}

	openIncompletePRs, err := getIncompletePRs(tx, params.TeamName, enum.PR_STATUS_OPEN)
	if err != nil {
		return entity.PullRequestGetIncompleteResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.PullRequestGetIncompleteResult{}, err
	}

	return entity.PullRequestGetIncompleteResult{
		TeamName:      params.TeamName,
		IncompletePRs: convertModelToEntity_ManyIncompletePRs(openIncompletePRs),
		FoundTeam:     true,
	}, nil
}

func getIncompletePRs(tx *sql.Tx, teamName string, status string) ([]model.PullRequestIncomplete, error) {
	incompletePRIds, err := getIncompletePRIds(tx, teamName, status)
	if err != nil {
		return []model.PullRequestIncomplete{}, err
	}

	incompletePRsInfo, err := getIncompletePRsMainInfo(tx, incompletePRIds)
	if err != nil {
		return []model.PullRequestIncomplete{}, err
	}

	for idx := range incompletePRsInfo {
		prReviewers, err := getPRReviewers(tx, incompletePRsInfo[idx].PullRequestId)
		if err != nil {
			return []model.PullRequestIncomplete{}, err
		}
		incompletePRsInfo[idx].AssignedReviewers = prReviewers
	}
	return incompletePRsInfo, nil
}

func getIncompletePRIds(tx *sql.Tx, teamName string, status string) ([]string, error) {
	query := `	SELECT
					pr.pull_request_id
				FROM
					pull_requests pr
				LEFT JOIN pull_requests__m2m__users m2m
					ON pr.pull_request_id = m2m.pull_request_id
				LEFT JOIN users reviewers
					ON reviewers.user_id = m2m.reviewer_id
				WHERE pr.status = $1
				GROUP BY pr.pull_request_id
				HAVING COUNT(reviewers.user_id) FILTER (
					WHERE reviewers.is_active = true
						AND EXISTS (
							SELECT 1
							FROM users team_users
							WHERE team_users.team_name = $2
								AND reviewers.user_id = team_users.user_id
						)
				) < $3`
	rows, err := tx.Query(query, status, teamName, consts.MAX_REVIEWERS_PER_PULL_REQUEST)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	resultPRIds := make([]string, 0, 10)
	var currentId string
	for rows.Next() {
		err := rows.Scan(&currentId)
		if err != nil {
			return []string{}, err
		}
		resultPRIds = append(resultPRIds, currentId)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return resultPRIds, nil
}

func getIncompletePRsMainInfo(tx *sql.Tx, pullRequestIds []string) ([]model.PullRequestIncomplete, error) {
	query := `	SELECT 
					pr.pull_request_id,
					pr.pull_request_name,
					pr.author_id,
					pr.status,
					pr.created_at
				FROM pull_requests pr
				INNER JOIN UNNEST($1::text[]) AS pr_ids(pull_request_id)
					ON pr.pull_request_id = pr_ids.pull_request_id`

	rows, err := tx.Query(query, pq.Array(pullRequestIds))
	if err != nil {
		return []model.PullRequestIncomplete{}, err
	}
	defer rows.Close()

	resultPRInfos := make([]model.PullRequestIncomplete, 0, 10)
	var currentPRInfo model.PullRequestIncomplete
	for rows.Next() {
		err := rows.Scan(
			&currentPRInfo.PullRequestId,
			&currentPRInfo.PullRequestName,
			&currentPRInfo.AuthorId,
			&currentPRInfo.Status,
			&currentPRInfo.CreatedAt,
		)
		if err != nil {
			return []model.PullRequestIncomplete{}, err
		}
		resultPRInfos = append(resultPRInfos, currentPRInfo)
	}

	err = rows.Err()
	if err != nil {
		return []model.PullRequestIncomplete{}, err
	}
	return resultPRInfos, nil
}

func getPRReviewers(tx *sql.Tx, prId string) ([]model.TeamMember, error) {
	query := `	SELECT
					users.user_id,
					users.username,
					users.is_active
				FROM pull_requests pr
				INNER JOIN pull_requests__M2M__users m2m
					ON m2m.pull_request_id = pr.pull_request_id
				INNER JOIN users 
					ON users.user_id = m2m.reviewer_id
				WHERE pr.pull_request_id = $1`

	rows, err := tx.Query(query, prId)
	if err != nil {
		return []model.TeamMember{}, err
	}
	defer rows.Close()

	resultReviewers := make([]model.TeamMember, 0, consts.MAX_REVIEWERS_PER_PULL_REQUEST)
	var currentReviewer model.TeamMember
	for rows.Next() {
		err := rows.Scan(
			&currentReviewer.UserId,
			&currentReviewer.Username,
			&currentReviewer.IsActive,
		)
		if err != nil {
			return []model.TeamMember{}, err
		}
		resultReviewers = append(resultReviewers, currentReviewer)
	}

	err = rows.Err()
	if err != nil {
		return []model.TeamMember{}, err
	}
	return resultReviewers, nil
}

func convertModelToEntity_ManyIncompletePRs(prs []model.PullRequestIncomplete) []entity.PullRequestIncomplete {
	result := make([]entity.PullRequestIncomplete, 0, len(prs))
	for _, pr := range prs {
		result = append(result, convertModelToEntity_OneIncompletePR(pr))
	}
	return result
}

func convertModelToEntity_OneIncompletePR(pr model.PullRequestIncomplete) entity.PullRequestIncomplete {
	return entity.PullRequestIncomplete{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: convertModelToEntity_ManyTeamMembers(pr.AssignedReviewers),
		CreatedAt:         pr.CreatedAt.Format(consts.FORMAT_LAYOUT_DATE_TIME),
	}
}
