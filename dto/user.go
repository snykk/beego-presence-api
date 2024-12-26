package dto

import (
	"time"

	"github.com/snykk/beego-presence-api/models"
)

type RegisterRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Department int    `json:"department"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DepartmentId int       `json:"department_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func FromUserModelToRegisterResponse(u models.User) RegisterResponse {
	return RegisterResponse{
		Id:           u.Id,
		Name:         u.Name,
		Email:        u.Email,
		DepartmentId: u.Department.Id,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
