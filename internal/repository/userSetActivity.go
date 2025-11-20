package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
)

func (repo *repository) UserSetActivity(ctx context.Context, params entity.UserSetActivityParams) (entity.UserSetActivityResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.UserSetActivityResult{}, nil
	}
	defer tx.Rollback()

	userExists, err := doesUserExists(tx, params.UserId)
	if err != nil {
		return entity.UserSetActivityResult{}, err
	}
	if !userExists {
		return entity.UserSetActivityResult{Found: false}, nil
	}

	err = updateUserActivity(tx, params)
	if err != nil {
		return entity.UserSetActivityResult{}, err
	}

	resultUser, err := getUserById(tx, params.UserId)
	if err != nil {
		return entity.UserSetActivityResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.UserSetActivityResult{}, nil
	}

	return entity.UserSetActivityResult{
		User:  convertModelToEntity_User(resultUser),
		Found: true,
	}, nil
}

func doesUserExists(tx *sql.Tx, userId string) (bool, error) {
	query := `SELECT user_id
				FROM users
				WHERE user_id = $1`
	var resId string
	err := tx.QueryRow(query, userId).Scan(&resId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func updateUserActivity(tx *sql.Tx, params entity.UserSetActivityParams) error {
	query := `UPDATE users
				SET is_active = $1
				WHERE user_id = $2`
	_, err := tx.Exec(query, params.NewActivenessValue, params.UserId)
	return err
}

func getUserById(tx *sql.Tx, userId string) (model.User, error) {
	query := `SELECT user_id, username, team_name, is_active
				FROM users
				WHERE user_id = $1`
	var user model.User
	err := tx.QueryRow(query, userId).Scan(
		&user.UserId,
		&user.Username,
		&user.TeamName,
		&user.IsActive,
	)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func convertModelToEntity_User(user model.User) entity.User {
	return entity.User{
		UserId:   user.UserId,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
