package handlers

import (
	"context"
	"net/http"
	"pr-reviewer-assignment-service/internal/dto"
	"pr-reviewer-assignment-service/internal/entity"
	"pr-reviewer-assignment-service/internal/enum"

	"github.com/labstack/echo/v4"
)

func (h *usersHandler) SetIsActiveProperty(c echo.Context) error {
	var reqDTO dto.UsersSetIsActive_Request
	if err := c.Bind(&reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}
	if err := c.Validate(reqDTO); err != nil {
		return c.JSON(http.StatusBadRequest,
			dto.NewErrorResponse(enum.ERROR_CODE_BAD_REQUEST, err.Error()))
	}

	updateParams := convertDTOToEntity_UserSetActivity(reqDTO)
	userUpdateResult, err := h.repo.UserSetActivity(context.TODO(), updateParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			dto.NewErrorResponse(enum.ERROR_CODE_INTERNAL_SERVER_ERROR, err.Error()))
	}
	if !userUpdateResult.Found {
		return c.JSON(http.StatusNotFound,
			dto.NewErrorResponse(enum.ERROR_CODE_NOT_FOUND, "Пользователь не найден"))
	}

	return c.JSON(http.StatusOK, convertEntityToDTO_User(userUpdateResult.User))
}

func convertDTOToEntity_UserSetActivity(reqDTO dto.UsersSetIsActive_Request) entity.UserSetActivityParams {
	return entity.UserSetActivityParams{
		UserId:             reqDTO.UserId,
		NewActivenessValue: *reqDTO.NewActivenessValue,
	}
}

func convertEntityToDTO_User(user entity.User) dto.User_Response {
	return dto.User_Response{
		UserId:   user.UserId,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}
