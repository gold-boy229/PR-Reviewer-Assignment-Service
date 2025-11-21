package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *pullRequestHandler) ReassignPullRequest(c echo.Context) error {
	var reqDTO dto.PullRequestReassign_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	reassignParams := convertDTOToEntity_PRReassign(reqDTO)
	reassignResult, err := h.repo.PullRequestReassign(context.TODO(), reassignParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !reassignResult.FoundPR {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "PR не найден"))
	}
	if !reassignResult.FoundOldReviewer {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "Пользователь не найден"))
	}
	if reassignResult.IsPullRequestMerged {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_PR_MERGED, "Нельзя менять после MERGED"))
	}
	if !reassignResult.IsOldReviewerAssigned {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_ASSIGNED, "Пользователь не был назначен ревьювером"))
	}
	if !reassignResult.FoundCandidate {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_NO_CANDIDATE, "Нет доступных кандидатов"))
	}

	return c.JSON(
		http.StatusOK,
		dto.PullRequestReassign_Response{
			PullRequest_Response: convertEntityToDTO_PullRequest(reassignResult.PullRequest),
			NewReviewerId:        reassignResult.NewReviewerId,
		},
	)
}

func convertDTOToEntity_PRReassign(reqDTO dto.PullRequestReassign_Request) entity.PullRequestReassignParams {
	return entity.PullRequestReassignParams{
		PullRequestId: reqDTO.PullRequestId,
		OldReviewerId: reqDTO.OldReviewerId,
	}
}
