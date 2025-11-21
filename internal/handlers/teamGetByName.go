package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *teamHandler) GetTeamByName(c echo.Context) error {
	var reqDTO dto.TeamNameQuery_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	teamSearchParams := convertDTOToEntity_TeamSearch(reqDTO)
	teamSearchResult, err := h.repo.GetTeamByName(context.TODO(), teamSearchParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !teamSearchResult.Found {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "Команда не найдена"))
	}

	return c.JSON(
		http.StatusOK,
		dto.TeamGet_Response(convertEntityToDTO_Team(teamSearchResult.Team)),
	)
}

func convertDTOToEntity_TeamSearch(reqDTO dto.TeamNameQuery_Request) entity.TeamSearchParams {
	return entity.TeamSearchParams{
		TeamName: reqDTO.TeamName,
	}
}
