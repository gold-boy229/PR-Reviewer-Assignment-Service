package handlers

import (
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
)

func convertEntityToDTO_Team(team entity.Team) dto.Team_Response {
	return dto.Team_Response{
		TeamName: team.TeamName,
		Members:  convertEntityToDTO_ManyTeamMembers(team.Members),
	}
}

func convertEntityToDTO_ManyTeamMembers(teamMembers []entity.TeamMember) []dto.TeamMember_Response {
	resultMembers := make([]dto.TeamMember_Response, 0, len(teamMembers))
	for _, member := range teamMembers {
		resultMembers = append(resultMembers, convertEntityToDTO_OneTeamMember(member))
	}
	return resultMembers
}

func convertEntityToDTO_OneTeamMember(tm entity.TeamMember) dto.TeamMember_Response {
	return dto.TeamMember_Response{
		UserId:   tm.UserId,
		Username: tm.Username,
		IsActive: tm.IsActive,
	}
}
