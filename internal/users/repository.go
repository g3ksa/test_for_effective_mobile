package users

import "UserService/internal/models"

type GetParams struct {
	Email    string
	Password string
}

type Repository interface {
	Create(user models.User) (*models.User, error)
	Edit(req EditUserRequest) (*models.User, error)
	Delete(userId int) error
	GetWithFilter(filters map[string]string, page, pageSize int) (*[]models.User, error)
}
