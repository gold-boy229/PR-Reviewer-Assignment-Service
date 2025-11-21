package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"
)

type usersProvider interface {
	GetUserAssignedPullRequests(context.Context, entity.UserGetAssignedPRParams) (entity.UserGetAssignedPRResult, error)
	UserSetActivity(context.Context, entity.UserSetActivityParams) (entity.UserSetActivityResult, error)
}

type usersHandler struct {
	repo usersProvider
}

func NewUsersHandler(repo usersProvider) *usersHandler {
	return &usersHandler{repo: repo}
}
