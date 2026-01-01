package main

import (
	"log/slog"
	"os"

	"github.com/muflihunaf/boilerplate-go/internal/app"

	_ "github.com/muflihunaf/boilerplate-go/docs" // Swagger docs
)

// @title           Boilerplate Go API
// @version         1.0
// @description     A production-ready Go boilerplate with JWT authentication.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/muflihunaf/boilerplate-go
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
