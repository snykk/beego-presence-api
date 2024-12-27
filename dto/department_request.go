package dto

import "github.com/snykk/beego-presence-api/models"

type DepartmentRequest struct {
	Name string `json:"name" validate:"required"`
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
