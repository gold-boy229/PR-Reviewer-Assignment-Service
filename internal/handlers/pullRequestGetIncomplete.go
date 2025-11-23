package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *pullRequestHandler) GetIncompletePRs(c echo.Context) error {
	var reqDTO dto.PullRequestGetIncomple_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	params := convertDTOToEntity_PRGetIncomplete(reqDTO)
	result, err := h.repo.PullRequestGetOpenIncompletePRs(context.TODO(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !result.FoundTeam {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "Команда не найдена"))
	}

	return c.JSON(http.StatusOK,
		dto.PullRequestGetIncomplete_Response{
			TeamName:      result.TeamName,
			IncompletePRs: convertEntityToDTO_ManyIncompletePRs(result.IncompletePRs),
		},
	)
}

func convertDTOToEntity_PRGetIncomplete(reqDTO dto.PullRequestGetIncomple_Request) entity.PullRequestGetIncompleteParams {
	return entity.PullRequestGetIncompleteParams{
		TeamName: reqDTO.TeamName,
	}
}

func convertEntityToDTO_ManyIncompletePRs(prs []entity.PullRequestIncomplete) []dto.PullRequestIncomplete_Response {
	result := make([]dto.PullRequestIncomplete_Response, 0, len(prs))
	for _, pr := range prs {
		result = append(result, convertEntityToDTO_OneIncompletePR(pr))
	}
	return result
}

func convertEntityToDTO_OneIncompletePR(pr entity.PullRequestIncomplete) dto.PullRequestIncomplete_Response {
	return dto.PullRequestIncomplete_Response{
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: convertEntityToDTO_ManyTeamMembers(pr.AssignedReviewers),
		CreatedAt:         pr.CreatedAt,
	}
}
