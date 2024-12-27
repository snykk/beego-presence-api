package dto

import (
	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/models"
)

type RegisterRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=50"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,securepwd"`
	DepartmentId int    `json:"department_id" validate:"required,min=1"`
}

func (r RegisterRequest) ToUserModel(md *models.Department) *models.User {
	return &models.User{
		Name:       r.Name,
		Email:      r.Email,
		Password:   r.Password,
		Department: md,
		Role:       constants.RoleEmployee, // default registered user is EMPLOYEE
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=50"`
	Email        string `json:"email" validate:"required,email"`
	DepartmentId int    `json:"department_id" validate:"required,min=1"`
}

func (u UserRequest) ToUserModel(mu *models.User, md *models.Department) *models.User {
	mu.Name = u.Name
	mu.Email = u.Email
	mu.Department = md
	return mu
}
