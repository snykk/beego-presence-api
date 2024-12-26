package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Schedule struct {
	Id         int         `orm:"auto" json:"id"`
	Name       string      `orm:"size(100)" json:"name"`
	Department *Department `orm:"rel(fk);column(department_id)" json:"department"` // ForeignKey to Department
	InTime     string      `orm:"size(8)" json:"in_time"`
	OutTime    string      `orm:"size(8)" json:"out_time"`
	Presences  []*Presence `orm:"reverse(many)" json:"presences"` // Reverse relationship with Presence
	Users      []*User     `orm:"reverse(many)" json:"users"`     // Reverse relationship with User
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)" json:"updated_at"`
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

	_, err = o.LoadRelated(&schedule, "Department")
	if err != nil {
		return nil, err
	}

	if isIncludePresenceList {
		_, err = o.LoadRelated(&schedule, "Presences")
		if err != nil {
			return nil, err
		}
	}

	if isIncludeUserList {
		_, err = o.LoadRelated(&schedule, "Presences")
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

func DeleteSchedule(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Schedule{Id: id})
	return err
}
