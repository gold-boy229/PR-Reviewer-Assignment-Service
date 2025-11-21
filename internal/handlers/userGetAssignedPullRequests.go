package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *usersHandler) GetUserAssignedPullRequests(c echo.Context) error {
	var reqDTO dto.UserGetReview_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	params := convertDTOToEntity_GetUserAssignedPRs(reqDTO)
	result, err := h.repo.GetUserAssignedPullRequests(context.TODO(), params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}

	return c.JSON(
		http.StatusOK,
		dto.UserGetReview_Response{
			UserId:       result.UserId,
			PullRequests: convertEntityToDTO_ManyShortPullRequests(result.PullRequests),
		},
	)
}

func convertDTOToEntity_GetUserAssignedPRs(reqDTO dto.UserGetReview_Request) entity.UserGetAssignedPRParams {
	return entity.UserGetAssignedPRParams{
		UserId: reqDTO.UserId,
	}
}

func convertEntityToDTO_ManyShortPullRequests(shortPRs []entity.PullRequestShort) []dto.PullRequestShort_Response {
	result := make([]dto.PullRequestShort_Response, 0, len(shortPRs))
	for _, pr := range shortPRs {
		result = append(result, convertEntityToDTO_OneShortPullRequest(pr))
	}
	return result
}

func convertEntityToDTO_OneShortPullRequest(pr entity.PullRequestShort) dto.PullRequestShort_Response {
	return dto.PullRequestShort_Response{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
	}
}
