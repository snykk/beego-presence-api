package dto

import "github.com/snykk/beego-presence-api/models"

// DepartmentRequest represents the structure of a department request
// @Description DepartmentRequest represents the structure of a department request
type DepartmentRequest struct {
	Name string `json:"name" validate:"required" example:"Human Resources"` // Department name
}

func (d *DepartmentRequest) ToDepartmentModel() *models.Department {
	return &models.Department{
		Name: d.Name,
	}
}

func (d *DepartmentRequest) ToDepartmentModelWithValue(md *models.Department) *models.Department {
	md.Name = d.Name
	return md
}
