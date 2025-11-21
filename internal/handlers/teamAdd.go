package handlers

import (
	"context"
	"fmt"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *teamHandler) AddTeam(c echo.Context) error {
	var reqDTO dto.TeamAdd_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest, dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	team := convertDTOToEntity_TeamAdd(reqDTO)
	teamSearchResult, err := h.repo.AddTeam(context.TODO(), team)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if teamSearchResult.Found {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_TEAM_EXISTS, fmt.Sprintf("team %q already exists", team.TeamName)))
	}

	return c.JSON(
		http.StatusCreated,
		dto.TeamAdd_Response{
			Team: convertEntityToDTO_Team(teamSearchResult.Team),
		},
	)
}

func convertDTOToEntity_TeamAdd(req dto.TeamAdd_Request) entity.Team {
	return entity.Team{
		TeamName: req.TeamName,
		Members:  convertDTOToEntity_ManyTeamMembers(req.Members),
	}
}

func convertDTOToEntity_ManyTeamMembers(teamMembers []dto.TeamMember_Request) []entity.TeamMember {
	resultMembers := make([]entity.TeamMember, 0, len(teamMembers))
	for _, member := range teamMembers {
		resultMembers = append(resultMembers, convertDTOToEntity_OneTeamMember(member))
	}
	return resultMembers
}

func convertDTOToEntity_OneTeamMember(tm dto.TeamMember_Request) entity.TeamMember {
	return entity.TeamMember{
		UserId:   tm.UserId,
		Username: tm.Username,
		IsActive: tm.IsActive,
	}
}
