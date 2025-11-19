package app

import (
	"fmt"
	"net/http"
	"pr-reviewer-assignment-service/internal/config"
	"pr-reviewer-assignment-service/internal/handlers"
	"pr-reviewer-assignment-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type App struct {
	echo *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Return the error so Echo's error handler can process it
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewApp() *App {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New(validator.WithRequiredStructEnabled())}
	return &App{echo: e}
}

func (app *App) Run() {
	configDB, err := config.ReadConfigDB()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Config loaded successfully\n ConfigDB = %+v\n", configDB)

	repo, err := repository.New(configDB)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Repository created successfully")

	var (
		teamHandler        teamHandler        = handlers.NewTeamHandler(repo)
		usersHandler       usersHandler       = handlers.NewUsersHandler(repo)
		pullRequestHandler pullRequestHandler = handlers.NewPullRequestHandler(repo)
	)

	app.echo.POST("/team/add", teamHandler.CreateTeam)
	app.echo.GET("/team/get", teamHandler.GetTeamByName)

	app.echo.POST("/users/setIsActive", usersHandler.SetIsActiveProperty)
	app.echo.GET("/users/getReview", usersHandler.GetUserAssignedPullRequests)

	app.echo.POST("/pullRequest/create", pullRequestHandler.CreatePullRequest)
	app.echo.POST("/pullRequest/merge", pullRequestHandler.MergePullRequest)
	app.echo.POST("/pullRequest/reassign", pullRequestHandler.ReassignPullRequest)

	app.echo.Start(":8080")
}
