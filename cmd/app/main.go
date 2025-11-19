package main

import "pr-reviewer-assignment-service/internal/app"

func main() {
	app := app.NewApp()
	app.Run()
}
