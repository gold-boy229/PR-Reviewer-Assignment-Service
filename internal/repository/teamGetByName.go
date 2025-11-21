package repository

import (
	"context"
	"database/sql"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
)

func (repo *repository) GetTeamByName(ctx context.Context, params entity.TeamSearchParams) (entity.TeamSearchResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}
	defer tx.Rollback()

	teamExists, err := doesTeamExist(tx, params.TeamName)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}
	if !teamExists {
		return entity.TeamSearchResult{FoundTeam: false}, nil
	}

	teamMembers, err := getTeamMembers(tx, params.TeamName)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	return entity.TeamSearchResult{
		Team: entity.Team{
			TeamName: params.TeamName,
			Members:  convertModelToEntity_ManyTeamMembers(teamMembers),
		},
		FoundTeam: true,
	}, nil
}

func getTeamMembers(tx *sql.Tx, teamName string) ([]model.TeamMember, error) {
	// prepare a slice of the required capacity
	rowsNumber, err := getTeamMembersNumber(tx, teamName)
	if err != nil {
		return []model.TeamMember{}, err
	}
	resultTeamMembers := make([]model.TeamMember, 0, rowsNumber)

	query := `SELECT user_id, username, is_active
				FROM users
				WHERE team_name = $1`

	rows, err := tx.Query(query, teamName)
	if err != nil {
		return []model.TeamMember{}, err
	}
	defer rows.Close()

	currentTM := model.TeamMember{}
	for rows.Next() {
		err := rows.Scan(&currentTM.UserId, &currentTM.Username, &currentTM.IsActive)
		if err != nil {
			return []model.TeamMember{}, err
		}
		resultTeamMembers = append(resultTeamMembers, currentTM)
	}

	err = rows.Err()
	if err != nil {
		return []model.TeamMember{}, err
	}
	return resultTeamMembers, nil
}

func getTeamMembersNumber(tx *sql.Tx, teamName string) (int, error) {
	query := `SELECT COUNT(*)
				FROM users
				WHERE team_name = $1`
	var rowsNumber int
	err := tx.QueryRow(query, teamName).Scan(&rowsNumber)
	if err != nil {
		return 0, err
	}
	return rowsNumber, nil
}

func convertModelToEntity_ManyTeamMembers(teamMembers []model.TeamMember) []entity.TeamMember {
	resultMembers := make([]entity.TeamMember, 0, len(teamMembers))
	for _, member := range teamMembers {
		resultMembers = append(resultMembers, convertModelToEntity_OneTeamMember(member))
	}
	return resultMembers
}

func convertModelToEntity_OneTeamMember(tm model.TeamMember) entity.TeamMember {
	return entity.TeamMember{
		UserId:   tm.UserId,
		Username: tm.Username,
		IsActive: tm.IsActive,
	}
}
