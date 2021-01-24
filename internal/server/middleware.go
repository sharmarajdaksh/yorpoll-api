package server

import (
	"time"

	"github.com/sharmarajdaksh/yorpoll-api/internal/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

var corsMiddleware = cors.New(cors.Config{
	AllowOrigins: "*",
	AllowHeaders: "Origin, Content-Type, Accept",
	AllowMethods: "GET, POST, PUT, DELETE",
})

var limiterMiddleware = limiter.New(limiter.Config{
	Max:        2,
	Expiration: 1 * time.Second,
	KeyGenerator: func(c *fiber.Ctx) string {
		return c.IP()
	},
	LimitReached: func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusTooManyRequests)
	},
})

func loggerMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	log.Logger.Info().Str("method", c.Method()).Str("path", c.Path()).Str("protocol", c.Protocol()).Str("originalUrl", c.OriginalURL()).Str("ip", c.IP()).Str("hostname", c.Hostname()).Int("statusCode", c.Response().StatusCode()).Str("requestID", c.Get(fiber.HeaderXRequestID)).Msg("received request")

	return err
}

func headersMiddleware(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.Next()
}

var requestIDMiddleware = requestid.New()
