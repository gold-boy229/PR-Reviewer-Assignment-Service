package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer-assignment-service/internal/consts"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
	"strings"
)

func (repo *repository) PullRequestCreate(ctx context.Context, params entity.PullRequestCreateParams) (entity.PullRequestCreateResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}
	defer tx.Rollback()

	teamName, err := getAuthorsTeamName(tx, params.AuthorId)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}
	if teamName == "" {
		return entity.PullRequestCreateResult{FoundAuthorAndTeam: false}, nil
	}

	prExiests, err := doesPullRequestExist(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}
	if prExiests {
		return entity.PullRequestCreateResult{FoundAuthorAndTeam: true, FoundPR: true}, nil
	}

	err = insertPullRequest(tx, params)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}

	prReviewersIds, err := findPRReviewersIds(tx, teamName, params.AuthorId, consts.MAX_REVIEWERS_PER_PULL_REQUEST)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}

	err = assignReviewersToPullRequest(tx, params.PullRequestId, prReviewersIds)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}

	resultPR, err := getPullRequestById(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.PullRequestCreateResult{}, err
	}

	return entity.PullRequestCreateResult{
		PullRequest:        convertModelToEntity_PullRequest(resultPR),
		FoundAuthorAndTeam: true,
		FoundPR:            false,
	}, nil
}

func getAuthorsTeamName(tx *sql.Tx, authorId string) (string, error) {
	query := `SELECT team_name
				FROM users
				WHERE user_id = $1`
	var teamName string
	err := tx.QueryRow(query, authorId).Scan(&teamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return teamName, nil
}

func doesPullRequestExist(tx *sql.Tx, prId string) (bool, error) {
	query := `SELECT pull_request_id
				FROM pull_requests
				WHERE pull_request_id = $1`
	var resPRId string
	err := tx.QueryRow(query, prId).Scan(&resPRId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func insertPullRequest(tx *sql.Tx, params entity.PullRequestCreateParams) error {
	query := `INSERT INTO pull_requests(
				pull_request_id,
				pull_request_name,
				author_id)
				VALUES ($1, $2, $3)`
	_, err := tx.Exec(
		query,
		params.PullRequestId,
		params.PullRequestName,
		params.AuthorId,
	)
	return err
}

func findPRReviewersIds(tx *sql.Tx, teamName string, authorId string, maxReviewersNum int) ([]string, error) {
	query := `SELECT user_id
				FROM users
				WHERE team_name = $1
					AND is_active = true
					AND user_id != $2
				LIMIT $3`
	rows, err := tx.Query(query, teamName, authorId, maxReviewersNum)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	resReviewersIds := make([]string, 0, maxReviewersNum)
	var currentUserId string
	for rows.Next() {
		err = rows.Scan(&currentUserId)
		if err != nil {
			return []string{}, err
		}
		resReviewersIds = append(resReviewersIds, currentUserId)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return resReviewersIds, nil
}

func assignReviewersToPullRequest(tx *sql.Tx, prId string, reviewersIds []string) error {
	query, args := prepareQueryAndArgs_assignReviewersToPR(prId, reviewersIds)
	_, err := tx.Exec(query, args...)
	return err
}

func prepareQueryAndArgs_assignReviewersToPR(prId string, reviewersIds []string) (string, []interface{}) {
	singleRowArgs := [...]interface{}{"pull_request_id", "reviewer_id"}
	const N = len(singleRowArgs)
	args := make([]interface{}, 0, len(reviewersIds)*N)

	if len(reviewersIds) == 0 {
		return "", args
	}

	var sb strings.Builder
	sb.WriteString(`INSERT INTO pull_requests__M2M__users (pull_request_id, reviewer_id) VALUES `)
	for i, reviewerId := range reviewersIds {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("($%d, $%d)", i*N+1, i*N+2))
		args = append(args, prId, reviewerId)
	}
	return sb.String(), args
}

func getPullRequestById(tx *sql.Tx, prId string) (model.PullRequest, error) {
	query := `	SELECT 
					pull_request_id,
					pull_request_name,
					author_id,
					status,
					created_at,
					merged_at
				FROM pull_requests
				WHERE pull_request_id = $1`
	var resPR model.PullRequest
	err := tx.QueryRow(query, prId).Scan(
		&resPR.PullRequestId,
		&resPR.PullRequestName,
		&resPR.AuthorId,
		&resPR.Status,
		&resPR.CreatedAt,
		&resPR.MergedAt,
	)
	if err != nil {
		return model.PullRequest{}, err
	}

	assignedReviewers, err := getAssignedReviewersIds(tx, prId)
	if err != nil {
		return model.PullRequest{}, err
	}

	resPR.AssignedReviewers = assignedReviewers
	return resPR, nil
}

func getAssignedReviewersIds(tx *sql.Tx, prId string) ([]string, error) {
	query := `	SELECT reviewer_id
				FROM pull_requests__M2M__users
				WHERE pull_request_id = $1`
	rows, err := tx.Query(query, prId)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	resIds := make([]string, 0, consts.MAX_REVIEWERS_PER_PULL_REQUEST)
	var curId string
	for rows.Next() {
		err := rows.Scan(&curId)
		if err != nil {
			return []string{}, err
		}
		resIds = append(resIds, curId)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return resIds, nil
}

func convertModelToEntity_PullRequest(pr model.PullRequest) entity.PullRequest {
	resPR := entity.PullRequest{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt.Format(consts.FORMAT_LAYOUT_DATE_TIME),
		MergedAt:          "",
	}
	if pr.MergedAt.Valid {
		resPR.MergedAt = pr.MergedAt.Time.Format(consts.FORMAT_LAYOUT_DATE_TIME)
	}
	return resPR
}
