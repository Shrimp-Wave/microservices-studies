package main

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type PersistenceLayer interface {
	CreateUser(request NewUserRequest) (*UserModel, error)
	GetById(id string) (*UserModel, error)
}

type InMemoryPersistenceLayer struct {
	data  map[string]*UserModel
	mutex sync.Mutex
}

func NewInMemoryPersistenceLayer() *InMemoryPersistenceLayer {
	data := make(map[string]*UserModel)
	data["satan"] = &UserModel{
		Id:       "satan",
		Email:    "stan@gmail.com.666",
		Username: "Stan Devile",
	}
	return &InMemoryPersistenceLayer{
		data: data,
	}
}

func (p *InMemoryPersistenceLayer) CreateUser(request NewUserRequest) (*UserModel, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id, _ := uuid.NewUUID()

	p.data[id.String()] = NewUserModel(
		id.String(),
		request.Username,
		request.Email,
		request.Password,
	)

	return p.data[id.String()], nil
}

func (p *InMemoryPersistenceLayer) GetById(id string) (user *UserModel, err error) {
	user, exists := p.data[id]
	if !exists {
		err = errors.New("user not found with id: " + id)
	}
	return
}
