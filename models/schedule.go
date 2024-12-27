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
	Users      []*User     `orm:"reverse(many)"` // Reverse relationship with User
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

// func init() {
// 	orm.RegisterModel(new(Schedule))
// }

func GetAllSchedules(isIncludePresenceList, isIncludeUserList bool) ([]*Schedule, error) {
	o := orm.NewOrm()
	var schedules []*Schedule
	_, err := o.QueryTable(new(Schedule)).RelatedSel("Department").All(&schedules)
	if err != nil {
		return nil, err
	}

	for i := range schedules {
		if isIncludePresenceList {
			_, err := o.LoadRelated(schedules[i], "Presences")
			if err != nil {
				return nil, err
			}
		}

		if isIncludeUserList {
			_, err := o.LoadRelated(schedules[i], "Users")
			if err != nil {
				return nil, err
			}
		}
	}
	return schedules, nil
}

func GetScheduleById(id int, isIncludePresenceList, isIncludeUserList bool) (*Schedule, error) {
	o := orm.NewOrm()
	schedule := &Schedule{Id: id}
	err := o.Read(schedule)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(schedule, "Department")
	if err != nil {
		return nil, err
	}

	if isIncludePresenceList {
		_, err = o.LoadRelated(schedule, "Presences")
		if err != nil {
			return nil, err
		}
	}

	if isIncludeUserList {
		_, err = o.LoadRelated(schedule, "Presences")
		if err != nil {
			return nil, err
		}
	}

	return schedule, nil
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

func DeleteSchedule(id int) (int64, error) {
	o := orm.NewOrm()
	affectedRows, err := o.Delete(&Schedule{Id: id})
	return affectedRows, err
}
