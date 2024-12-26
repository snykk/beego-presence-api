package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Department struct {
	Id        int         `orm:"auto" json:"id"`
	Name      string      `orm:"size(100)" json:"name"`
	Users     []*User     `orm:"reverse(many)" json:"users"`     // Reverse relationship with User
	Schedules []*Schedule `orm:"reverse(many)" json:"schedules"` // Reverse relationship with Schedule
	CreatedAt time.Time   `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time   `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// func init() {
// 	orm.RegisterModel(new(Department))
// }

func GetAllDepartments() ([]Department, error) {
	o := orm.NewOrm()
	var departments []Department
	_, err := o.QueryTable(new(Department)).All(&departments)
	return departments, err
}

func GetDepartmentById(id int) (Department, error) {
	o := orm.NewOrm()
	department := Department{Id: id}
	err := o.Read(&department)
	return department, err
}

func CreateDepartment(department *Department) error {
	o := orm.NewOrm()
	_, err := o.Insert(department)
	return err
}

func UpdateDepartment(department *Department) error {
	o := orm.NewOrm()
	_, err := o.Update(department)
	return err
}

func DeleteDepartment(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Department{Id: id})
	return err
}
