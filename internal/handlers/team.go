package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"
)

type teamProvider interface {
	AddTeam(context.Context, entity.Team) (entity.Team, bool, error)
	GetTeamByName(context.Context, entity.TeamSearchParams) (entity.TeamSearchResult, error)
}

type teamHandler struct {
	repo teamProvider
}

func NewTeamHandler(repo teamProvider) *teamHandler {
	return &teamHandler{repo: repo}
}
