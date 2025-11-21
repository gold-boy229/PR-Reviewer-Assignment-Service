package repository

import (
	"context"
	"database/sql"
	"fmt"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/model"
	"strings"

	"github.com/lib/pq"
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
		return entity.TeamSearchResult{FoundTeam: true}, nil
	}

	err = insertTeamName(tx, team.TeamName)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}

	duplicateTMIds, err := getTeamMembersDuplicatesInDB(tx, team.Members)
	if err != nil {
		return entity.TeamSearchResult{}, err
	}
	if len(duplicateTMIds) != 0 {
		return entity.TeamSearchResult{
			ConflictingUserIds: duplicateTMIds,
			FoundTeam:          false,
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
		Team:      team,
		FoundTeam: false,
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

func getTeamMembersDuplicatesInDB(tx *sql.Tx, teamMembers []entity.TeamMember) ([]string, error) {
	teamMemberIds := getTeamMemberIds(teamMembers)

	duplicateIds := make([]string, 0)
	var curUserId string

	query := `	SELECT u.user_id
				FROM users u
				INNER JOIN UNNEST($1::text[]) AS ids(user_id) 
					ON u.user_id = ids.user_id;`
	rows, err := tx.Query(query, pq.Array(teamMemberIds))
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&curUserId)
		if err != nil {
			return []string{}, err
		}
		duplicateIds = append(duplicateIds, curUserId)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return duplicateIds, nil
}

func getTeamMemberIds(teamMembers []entity.TeamMember) []string {
	result := make([]string, 0, len(teamMembers))
	for _, tm := range teamMembers {
		result = append(result, tm.UserId)
	}
	return result
}

func insertTeamMembers(tx *sql.Tx, team model.Team) error {
	query, args := prepareQueryAndArgs_insertTeamMembers(team)
	_, err := tx.Exec(query, args...)
	return err
}

func prepareQueryAndArgs_insertTeamMembers(team model.Team) (string, []interface{}) {
	singleRowArgs := [...]interface{}{"user_id", "username", "team_name", "is_active"}
	const N = len(singleRowArgs)
	args := make([]interface{}, 0, len(team.Members)*N)

	if len(team.Members) == 0 {
		return "", args
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
	return sb.String(), args
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
