package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (c CustomValidator) Validate(i interface{}) error {
	if err := c.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewCustomValidator() *CustomValidator {
	validate := validator.New()

	// Add Custom validations
	// EXAMPLE: validate.RegisterValidation("is-awesome", ValidateMyVal)

	return &CustomValidator{validate}
}
