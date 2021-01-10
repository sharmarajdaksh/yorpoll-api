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
	app.Use(headersMiddleware)
	app.Use(loggerMiddleware)

	healthCheckHandler := handler{dbc: d}
	app.Get("/healthcheck", healthCheckHandler.healthCheck)

	api := app.Group("/api")

	v1 := api.Group("/v1")
	registerV1Routes(v1, d)

	// initialize validator
	initializeValidator()

	return app
}

func registerV1Routes(r fiber.Router, d db.Connection) {
	v1Handler := handler{dbc: d}
	r.Get("/poll/:pollID", v1Handler.getPoll)
	r.Post("/poll", v1Handler.postPoll)
	r.Delete("/poll/:pollID", v1Handler.deletePoll)
	r.Put("/vote/:pollID/:optionID", v1Handler.putVote)
}
