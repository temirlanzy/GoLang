package repository

import (
	"practice3/internal/repository/_postgres"
	"practice3/internal/repository/_postgres/users"
	"practice3/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	GetUserByID(int) (*modules.User, error)
	CreateUser(modules.User) (int, error)
	UpdateUser(int, modules.User) error
	DeleteUser(int) error
}

type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
