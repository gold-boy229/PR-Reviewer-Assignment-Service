package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"

	"github.com/labstack/echo/v4"
)

type usersProvider interface {
	UserSetActivity(context.Context, entity.UserSetActivityParams) (entity.UserSetActivityResult, error)
}

type usersHandler struct {
	repo usersProvider
}

func NewUsersHandler(repo usersProvider) *usersHandler {
	return &usersHandler{repo: repo}
}

func (h *usersHandler) GetUserAssignedPullRequests(c echo.Context) error {
	return nil
}
