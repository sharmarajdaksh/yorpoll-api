package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sharmarajdaksh/yorpoll-api/config"
	"github.com/sharmarajdaksh/yorpoll-api/internal/db"
)

// Init initializes the server routes and handlers
func Init(c *config.Config, d db.Connection) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(corsMiddleware)
	app.Use(limiterMiddleware)
	app.Use(requestIDMiddleware)
	app.Use(headersMiddleware)
	app.Use(loggerMiddleware)

	healthCheckHandler := handler{dbc: d}
	app.Get("/healthcheck", healthCheckHandler.healthCheck)

	if c.Global.Env != config.Prod && c.Global.Env != config.Trace {
		app.Static("/swagger", "swaggerui")
	}

	api := app.Group("/api")

	v1 := api.Group("/v1")
	registerV1Routes(v1, d)

	// initialize validator
	initializeValidator()

	return app
}
