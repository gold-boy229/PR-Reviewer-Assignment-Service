package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
	"strings"
)

func (repo *repository) AddTeam(ctx context.Context, team entity.Team) (entity.TeamSearchResult, error) {
	tx, err := repo.Db.BeginTx(ctx, nil)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}
	defer tx.Rollback()

	teamExists, err := doesTeamExist(tx, team.TeamName)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}
	if teamExists {
		return entity.TeamSearchResult{Found: true}, nil
	}

	err = insertTeamName(tx, team.TeamName)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	if len(team.Members) == 0 {
		return entity.TeamSearchResult{
			Team:  team,
			Found: false,
		}, nil
	}
	err = insertTeamMembers(tx, convertEntityToModel_Team(team))
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	err = tx.Commit()
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	return entity.TeamSearchResult{
		Team:  team,
		Found: false,
	}, nil
}

func doesTeamExist(tx *sql.Tx, teamName string) (bool, error) {
	query := `SELECT team_name 
				FROM teams
				WHERE team_name = $1`
	var resTeamName string
	err := tx.QueryRow(query, teamName).Scan(&resTeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func insertTeamName(tx *sql.Tx, teamName string) error {
	query := `INSERT INTO teams (team_name) VALUES ($1)`
	_, err := tx.Exec(query, teamName)
	return err
}

func insertTeamMembers(tx *sql.Tx, team model.Team) error {
	query, args, err := prepareQueryAndArgs_insertTeamMembers(team)
	if err != nil {
		return err
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func prepareQueryAndArgs_insertTeamMembers(team model.Team) (string, []interface{}, error) {
	singleRowArgs := [...]interface{}{"user_id", "username", "team_name", "is_active"}
	const N = len(singleRowArgs)
	args := make([]interface{}, 0, len(team.Members)*N)

	if len(team.Members) == 0 {
		return "", args, nil
	}

	var sb strings.Builder
	sb.WriteString("INSERT INTO users (user_id, username, team_name, is_active) VALUES ")
	for i, tm := range team.Members {
		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", i*N+1, i*N+2, i*N+3, i*N+4))
		args = append(args, tm.UserId, tm.Username, team.TeamName, tm.IsActive)
	}
	return sb.String(), args, nil
}

func convertEntityToModel_Team(team entity.Team) model.Team {
	return model.Team{
		TeamName: team.TeamName,
		Members:  convertEntityToModel_ManyTeamMembers(team.Members),
	}
}

func convertEntityToModel_ManyTeamMembers(teamMembers []entity.TeamMember) []model.TeamMember {
	resultMembers := make([]model.TeamMember, 0, len(teamMembers))
	for _, member := range teamMembers {
		resultMembers = append(resultMembers, convertEntityToModel_OneTeamMember(member))
	}
	return resultMembers
}

func convertEntityToModel_OneTeamMember(tm entity.TeamMember) model.TeamMember {
	return model.TeamMember{
		UserId:   tm.UserId,
		Username: tm.Username,
		IsActive: tm.IsActive,
	}
}
