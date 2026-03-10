package interfaces

import (
	"hr-system/internal/models"

	"github.com/google/uuid"
)

type UserInterface interface {
	Register(user *models.User) error
	Login(email, password string) (map[string]interface{}, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uuid.UUID) (*models.User, error)
	UpdateUser(updates *models.User) (*models.User, error)
	DeactivateUser(id uuid.UUID) error
}
