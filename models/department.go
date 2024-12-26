package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Department struct {
	Id        int         `orm:"auto"`
	Name      string      `orm:"size(100)"`
	Users     []*User     `orm:"reverse(many)"` // Reverse relationship with User
	Schedules []*Schedule `orm:"reverse(many)"` // Reverse relationship with Schedule
	CreatedAt time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time   `orm:"auto_now;type(datetime)"`
}

// func init() {
// 	orm.RegisterModel(new(Department))
// }

func GetAllDepartments(isIncludeUserList, isIncludeScheduleList bool) ([]*Department, error) {
	o := orm.NewOrm()
	var departments []*Department
	// Fetch all departments
	_, err := o.QueryTable(new(Department)).All(&departments)
	if err != nil {
		return nil, err
	}

	// Load related Users and Schedules for each department
	for i := range departments {
		if isIncludeUserList {
			_, err = o.LoadRelated(departments[i], "Users")
			if err != nil {
				return nil, err
			}
		}

		if isIncludeScheduleList {
			_, err = o.LoadRelated(departments[i], "Schedules")
			if err != nil {
				return nil, err
			}
		}
	}

	return departments, nil
}

func GetDepartmentById(id int, isIncludeUserList, isIncludeScheduleList bool) (*Department, error) {
	o := orm.NewOrm()
	department := &Department{Id: id}
	err := o.Read(department)
	if err != nil {
		return nil, err
	}

	if isIncludeUserList {
		_, err = o.LoadRelated(department, "Users")
		if err != nil {
			return nil, err
		}
	}

	if isIncludeScheduleList {
		_, err = o.LoadRelated(department, "Schedules")
		if err != nil {
			return nil, err
		}
	}

	return department, nil
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
