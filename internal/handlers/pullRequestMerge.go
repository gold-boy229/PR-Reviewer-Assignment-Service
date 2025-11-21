package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *pullRequestHandler) MergePullRequest(c echo.Context) error {
	var reqDTO dto.PullRequestMerge_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	mergeParams := convertDTOToEntity_PRMerge(reqDTO)
	resultPRMerge, err := h.repo.PullRequestMerge(context.TODO(), mergeParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !resultPRMerge.FoundPR {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "PR не найден"))
	}

	return c.JSON(
		http.StatusOK,
		dto.PullRequestMerge_Response{
			PullRequest_Response: convertEntityToDTO_PullRequest(resultPRMerge.PullRequest),
		},
	)
}

func convertDTOToEntity_PRMerge(reqDTO dto.PullRequestMerge_Request) entity.PullRequestMergeParams {
	return entity.PullRequestMergeParams{
		PullRequestId: reqDTO.PullRequestId,
	}
}
