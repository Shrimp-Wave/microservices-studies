package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Service interface {
	CreateUser(c echo.Context) error
	GetUserById(c echo.Context) error
}

type UserService struct {
	database PersistenceLayer
	mail     *MailBroker
}

func NewUserService(layer PersistenceLayer, mail *MailBroker) *UserService {
	return &UserService{
		database: layer,
		mail:     mail,
	}
}

func (u *UserService) CreateUser(c echo.Context) error {
	requestBody := c.Get("validatedData")
	newUser, err := u.database.CreateUser(requestBody.(NewUserRequest))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	go u.mail.Send(NewMailPayload(newUser.Email, newUser.Username, "Registration Completed!"))

	c.Logger().Info("New user created with request body", "REQUEST_BODY", requestBody)
	return c.JSONPretty(http.StatusOK, newUser, "  ")
}

func (u *UserService) GetUserById(c echo.Context) error {
	userId := c.Param("id")
	user, err := u.database.GetById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSONPretty(http.StatusOK, user, "  ")
}
