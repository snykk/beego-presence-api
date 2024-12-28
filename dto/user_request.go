package dto

import (
	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/models"
)

// RegisterRequest represents the structure of a register request
// @Description RegisterRequest represents the structure of a register request
type RegisterRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=50" example:"Najib Fikri"`    // Name of the user
	Email        string `json:"email" validate:"required,email" example:"najibfikri@gmail.com"` // Email of the user
	Password     string `json:"password" validate:"required,securepwd" example:"Mys3cur3P@5s"`  // Password of the user
	DepartmentId int    `json:"department_id" validate:"required,min=1" example:"1"`            // ForeignKey to Department
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

// LoginRequest represents the structure of a login request
// @Description LoginRequest represents the structure of a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"najibfikri@gmail.com"` // Email of the user
	Password string `json:"password" validate:"required" example:"Mys3cur3P@5s"`            // Password of the user
}

// UserRequest represents the structure of a user request
// @Description UserRequest represents the structure of a user request
type UserRequest struct {
	Name         string `json:"name" validate:"required,min=3,max=50" example:"Najib Fikri"`    // Name of the user
	Email        string `json:"email" validate:"required,email" example:"najibfikri@gmail.com"` // Email of the user
	DepartmentId int    `json:"department_id" validate:"required,min=1" example:"1"`            // ForeignKey to Department
}

func (u UserRequest) ToUserModel(mu *models.User, md *models.Department) *models.User {
	mu.Name = u.Name
	mu.Email = u.Email
	mu.Department = md
	return mu
}
