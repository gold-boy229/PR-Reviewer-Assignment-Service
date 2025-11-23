package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *pullRequestHandler) AssignNewReviewers(c echo.Context) error {
	var reqDTO dto.PullRequestAssignReviewers_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	params := convertModelToEntity_PullRequestAssignNewReviewers(reqDTO)
	result, err := h.repo.PullRequestAssignReviewers(context.TODO(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !result.FoundPullRequest {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "PR не найден"))
	}
	if result.IsPullRequestMerged {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_PR_MERGED, "Нельзя добавлять ревьюеров после MERGED"))
	}
	if result.HasMaxReviewersAmount {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_PR_MAX_REVIEWERS, "У PR'а уже есть два ревьюера"))
	}
	if !result.FoundCandidate {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_NO_CANDIDATE, "Нет доступных кандидатов"))
	}

	return c.JSON(http.StatusOK,
		dto.PullRequestAssignReviewers_Response{
			PullRequestIncomplete: convertEntityToDTO_OneIncompletePR(result.PullRequestIncomplete),
		},
	)
}

func convertModelToEntity_PullRequestAssignNewReviewers(reqDTO dto.PullRequestAssignReviewers_Request) entity.PullRequestAssignReviewersParams {
	return entity.PullRequestAssignReviewersParams{
		PullRequestId: reqDTO.PullRequestId,
	}
}
