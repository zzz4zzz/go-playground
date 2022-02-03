package main

import (
	"github.com/zzz4zzz/go-playground/application"
	"github.com/zzz4zzz/go-playground/db"
)

func main() {
	var app = application.NewApplication(
		db.DatabaseConfiguration{
			Host:            "localhost",
			Port:            5432,
			User:            "user",
			Password:        "password",
			DatabaseName:    "playground",
			ApplicationName: "go-playground",
		}, "localhost:8080")

	err := app.Run()

	if err != nil {
		panic("Application could not be served")
	}
}
