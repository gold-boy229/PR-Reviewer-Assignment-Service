package handlers

import "github.com/labstack/echo/v4"

type teamProvider interface {
}

type teamHandler struct {
	repo teamProvider
}

func NewTeamHandler(repo teamProvider) *teamHandler {
	return &teamHandler{repo: repo}
}

func (h *teamHandler) CreateTeam(c echo.Context) error {
	return nil
}

func (h *teamHandler) GetTeamByName(c echo.Context) error {
	return nil
}
