package cmd

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/marmorag/supateam/docs"
	"github.com/marmorag/supateam/internal"
	"github.com/marmorag/supateam/internal/repository"
	"github.com/marmorag/supateam/internal/routes"
	"github.com/spf13/cobra"
	"log"
)

func executeServeCommand() {
	config := internal.GetConfig()

	if config.ApplicationEnvironment == "prod" {
		err := sentry.Init(sentry.ClientOptions{
			Environment: "prod",
			Dsn:         "https://90dc370dbd474422bdf2394c51e0d65e@o473284.ingest.sentry.io/6129025",
		})
		if err != nil {
			log.Fatalf("sentry.Init: %s", err)
		}
	}

	handlers := []routes.RouteHandler{
		routes.AuthRouteHandler{},
		routes.HealthcheckRouteHandler{},
		routes.UserRouteHandler{},
		routes.EventRouteHandler{},
		routes.ParticipationRouteHandler{},
		routes.TeamRouteHandler{},
	}
	app := fiber.New(fiber.Config{
		Prefork: config.ApplicationPrefork,
		AppName: config.ApplicationName,
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: config.CorsAllowOrigins,
		AllowHeaders: config.CorsAllowHeaders,
		AllowMethods: config.CorsAllowMethods,
	}))
	//app.Use(auth.NewAuthHeaderHandler())

	api := app.Group("/api", logger.New())
	api.Get("/swagger/*", swagger.Handler)

	for _, handler := range handlers {
		handler.Register(api)
	}

	defer repository.CloseConnection()

	log.Println("Serving supateam API : http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "HTTP REST API for supateam App.",
	Run: func(cmd *cobra.Command, args []string) {
		executeServeCommand()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
