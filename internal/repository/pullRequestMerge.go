package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"
)

func (repo *repository) PullRequestMerge(ctx context.Context, params entity.PullRequestMergeParams) (entity.PullRequestMergeResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.PullRequestMergeResult{}, err
	}
	defer tx.Rollback()

	pullRequestExists, err := doesPullRequestExist(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestMergeResult{}, err
	}
	if !pullRequestExists {
		return entity.PullRequestMergeResult{FoundPR: false}, nil
	}

	err = markPullRequestAsMerged(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestMergeResult{}, err
	}

	pullRequest, err := getPullRequestById(tx, params.PullRequestId)
	if err != nil {
		return entity.PullRequestMergeResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.PullRequestMergeResult{}, err
	}

	return entity.PullRequestMergeResult{
		PullRequest: convertModelToEntity_PullRequest(pullRequest),
		FoundPR:     true,
	}, nil
}

// Пометить PR как MERGED (идемпотентная операция)
func markPullRequestAsMerged(tx *sql.Tx, prId string) error {
	query := `	UPDATE pull_requests
				SET 
					status = $1,
					merged_at = COALESCE(merged_at, Now())
				WHERE pull_request_id = $2;`
	_, err := tx.Exec(query, enum.PR_STATUS_MERGED, prId)
	return err
}
