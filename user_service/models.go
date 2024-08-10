package main

import (
	"crypto/sha256"
	"time"
)

type NewUserRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=255"`
	Email    string `json:"email" validate:"required,email,gte=2,lte=255"`
	Password string `json:"password" validate:"required,gte=18,lte=18"`
}

type UserModel struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
	Salt     []byte `json:"-"`
}

func NewUserModel(id, username, email, password string) *UserModel {
	hash := sha256.New()
	hash.Write([]byte(password))
	salt := sha256.Sum256([]byte("Jolestest" + time.Now().String() + id))

	return &UserModel{
		Id:       id,
		Username: username,
		Email:    email,
		Password: hash.Sum(salt[:]),
		Salt:     salt[:],
	}
}
