package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *pullRequestHandler) CreatePullRequest(c echo.Context) error {
	var reqDTO dto.PullRequestCreate_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	prCreateParams := convertDTOToEntity_PRCreate(reqDTO)
	prCreateResult, err := h.repo.PullRequestCreate(context.TODO(), prCreateParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !prCreateResult.FoundAuthorAndTeam {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "Автор/команда не найдены"))
	}
	if prCreateResult.FoundPR {
		return c.JSON(http.StatusConflict,
			dto.NewErrorResponse(enum.ERROR_CODE_PR_EXISTS, "PR уже существует"))
	}

	return c.JSON(http.StatusCreated, convertEntityToDTO_PullRequest(prCreateResult.PullRequest))
}

func convertDTOToEntity_PRCreate(reqDTO dto.PullRequestCreate_Request) entity.PullRequestCreateParams {
	return entity.PullRequestCreateParams{
		PullRequestId:   reqDTO.PullRequestId,
		PullRequestName: reqDTO.PullRequestName,
		AuthorId:        reqDTO.AuthorId,
	}
}

func convertEntityToDTO_PullRequest(pr entity.PullRequest) dto.PullRequest_Response {
	return dto.PullRequest_Response{
		PullRequestId:        pr.PullRequestId,
		PullRequestName:      pr.PullRequestName,
		AuthorId:             pr.AuthorId,
		Status:               pr.Status,
		AssignedReviewersIds: pr.AssignedReviewers,
		CreatedAt:            pr.CreatedAt,
		MergedAt:             pr.MergedAt,
	}
}
