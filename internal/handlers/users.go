package handlers

import "github.com/labstack/echo/v4"

type usersProvider interface {
}

type usersHandler struct {
	repo usersProvider
}

func NewUsersHandler(repo usersProvider) *usersHandler {
	return &usersHandler{repo: repo}
}

func (h *usersHandler) SetIsActiveProperty(c echo.Context) error {
	return nil
}

func (h *usersHandler) GetUserAssignedPullRequests(c echo.Context) error {
	return nil
}
