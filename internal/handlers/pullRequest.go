package handlers

import (
	"context"
	"pr-reviewer-assignment-service/internal/entity"
)

type pullRequestProvider interface {
	PullRequestCreate(context.Context, entity.PullRequestCreateParams) (entity.PullRequestCreateResult, error)
	PullRequestMerge(context.Context, entity.PullRequestMergeParams) (entity.PullRequestMergeResult, error)
	PullRequestReassign(context.Context, entity.PullRequestReassignParams) (entity.PullRequestReassignResult, error)
}

type pullRequestHandler struct {
	repo pullRequestProvider
}

func NewPullRequestHandler(repo pullRequestProvider) *pullRequestHandler {
	return &pullRequestHandler{repo: repo}
}
