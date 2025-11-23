package model

import (
	"database/sql"
	"time"
)

type PullRequest struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          sql.NullTime
}

type PullRequestShort struct {
	PullRequestId   string
	PullRequestName string
	AuthorId        string
	Status          string
}

type PullRequestIncomplete struct {
	PullRequestId     string
	PullRequestName   string
	AuthorId          string
	Status            string
	AssignedReviewers []TeamMember
	CreatedAt         time.Time
}
