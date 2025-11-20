package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"

	"github.com/labstack/echo/v4"
)

type teamProvider interface {
	AddTeam(context.Context, entity.Team) (entity.Team, bool, error)
}

type teamHandler struct {
	repo teamProvider
}

func NewTeamHandler(repo teamProvider) *teamHandler {
	return &teamHandler{repo: repo}
}

func (h *teamHandler) GetTeamByName(c echo.Context) error {
	return nil
}
