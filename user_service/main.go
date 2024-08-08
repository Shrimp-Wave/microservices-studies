package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// TODO: Implement user creation send email via email service
// TODO: Implement user buy subscription calling RUST buying service
// TODO: Implement Python service to validate if transaction can be done successfully
// TODO: Implement database in everything

func main() {
	e := echo.New()

	e.GET("/", hello)
	e.Logger.Fatal(e.Start(":6723"))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
