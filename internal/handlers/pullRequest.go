package handlers

import "github.com/labstack/echo/v4"

type pullRequestProvider interface {
}

type pullRequestHandler struct {
	repo pullRequestProvider
}

func NewPullRequestHandler(repo pullRequestProvider) *pullRequestHandler {
	return &pullRequestHandler{repo: repo}
}

func (h *pullRequestHandler) CreatePullRequest(c echo.Context) error {
	return nil
}

func (h *pullRequestHandler) MergePullRequest(c echo.Context) error {
	return nil
}

func (h *pullRequestHandler) ReassignPullRequest(c echo.Context) error {
	return nil
}
