package main

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func RequestIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var requestId string

		if requestId = c.Request().Header.Get(requestIdHeader); strings.TrimSpace(requestId) == "" {
			requestUUID, _ := uuid.NewUUID()
			requestId = requestUUID.String()
		}

		c.Logger().Info("Adding new request id...", "REQUEST_ID=", requestId)
		c.Set(requestIdHeader, requestId)
		return next(c)
	}
}

func CreateUserValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger := c.Logger()

		var requestBody NewUserRequest

		if err := c.Bind(&requestBody); err != nil {
			logger.Error("Failed to bind user body", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(requestBody); err != nil {
			logger.Error("User validation failed", err)
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		logger.Info("Successfully validated user request")
		c.Set("validatedData", requestBody)
		return next(c)
	}
}
