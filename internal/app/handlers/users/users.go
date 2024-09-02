package users

import (
	"hamsterbot/internal/app/models"
)

type User interface {
	GetUserById(id int64) (map[string]interface{}, error)
	AddUser(user models.User) error
	UpdateUser(user models.User) error
	DeleteUser(id int64) error
}

type Endpoint struct {
	User User
}
