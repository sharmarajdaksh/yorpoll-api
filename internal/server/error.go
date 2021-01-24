package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sharmarajdaksh/yorpoll-api/internal/log"
)

var vldt *validator.Validate

func errorHandler(c *fiber.Ctx, err error) error {
	log.Logger.Error().Err(err).Str("requestid", c.Get(fiber.HeaderXRequestID)).Send()

	// Default 500 statuscode
	code := http.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}

	// Return statuscode with error message
	return c.Status(code).SendString(err.Error())
}

func validate(str interface{}) error {

	err := vldt.Struct(str)

	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fiber.ErrBadRequest
		}

		errs := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("Error in %s type field \"%s\". Passed value: \"%v\"", err.Kind(), err.Field(), err.Value()))
		}
		return fiber.NewError(http.StatusBadRequest, errs...)
	}

	return nil
}

func initializeValidator() {

	vldt = validator.New()

	vldt.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
}
