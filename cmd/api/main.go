package main

import (
	"log/slog"
	"os"

	"github.com/muflihunaf/boilerplate-go/internal/app"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	application, err := app.New()
	if err != nil {
		return err
	}

	return application.Run()
}
