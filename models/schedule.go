package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Schedule struct {
	Id         int         `orm:"auto"`
	Name       string      `orm:"size(100)"`
	Department *Department `orm:"rel(fk);column(department_id)"` // ForeignKey to Department
	InTime     string      `orm:"size(8)"`
	OutTime    string      `orm:"size(8)"`
	Presences  []*Presence `orm:"reverse(many)"` // Reverse relationship with Presence
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

// func init() {
// 	orm.RegisterModel(new(Schedule))
// }

func GetAllSchedules() ([]Schedule, error) {
	o := orm.NewOrm()
	var schedules []Schedule
	_, err := o.QueryTable(new(Schedule)).All(&schedules)
	return schedules, err
}

func GetScheduleById(id int) (Schedule, error) {
	o := orm.NewOrm()
	schedule := Schedule{Id: id}
	err := o.Read(&schedule)
	return schedule, err
}

func CreateSchedule(schedule *Schedule) error {
	o := orm.NewOrm()
	_, err := o.Insert(schedule)
	return err
}

func UpdateSchedule(schedule *Schedule) error {
	o := orm.NewOrm()
	_, err := o.Update(schedule)
	return err
}

func DeleteSchedule(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Schedule{Id: id})
	return err
}
