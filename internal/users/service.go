package users

import "UserService/internal/models"

type EditUserRequest struct {
	Id       int         `json:"id"`
	Field    string      `json:"field"`
	NewValue interface{} `json:"new_value"`
}

type Service interface {
	GetWithFilter(filters map[string]string, page, pageSize int) (*[]models.User, error)
	Delete(userId int) error
	Edit(EditUserRequest) (*models.User, error)
	Create(user models.User) (*models.User, error)
}
