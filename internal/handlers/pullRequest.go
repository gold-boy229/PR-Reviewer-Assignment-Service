package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"

	"github.com/labstack/echo/v4"
)

type pullRequestProvider interface {
	PullRequestCreate(context.Context, entity.PullRequestCreateParams) (entity.PullRequestCreateResult, error)
}

type pullRequestHandler struct {
	repo pullRequestProvider
}

func NewPullRequestHandler(repo pullRequestProvider) *pullRequestHandler {
	return &pullRequestHandler{repo: repo}
}

func (h *pullRequestHandler) MergePullRequest(c echo.Context) error {
	return nil
}

func (h *pullRequestHandler) ReassignPullRequest(c echo.Context) error {
	return nil
}
