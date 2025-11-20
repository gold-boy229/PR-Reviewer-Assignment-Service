package app

import "github.com/labstack/echo/v4"

type pullRequestHandler interface {
	CreatePullRequest(c echo.Context) error
	MergePullRequest(c echo.Context) error
	ReassignPullRequest(c echo.Context) error
}

type teamHandler interface {
	AddTeam(c echo.Context) error
	GetTeamByName(c echo.Context) error
}

type usersHandler interface {
	SetIsActiveProperty(c echo.Context) error
	GetUserAssignedPullRequests(c echo.Context) error
}
