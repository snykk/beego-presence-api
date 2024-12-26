package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Presence represents the presence table in the database
type Presence struct {
	Id        int       `orm:"auto"`
	User      *User     `orm:"rel(fk)"` // ForeignKey to User
	Schedule  *Schedule `orm:"rel(fk)"` // ForeignKey to Schedule
	Type      string    `orm:"size(10)"`
	Status    string    `orm:"size(50)"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

// func init() {
// 	orm.RegisterModel(new(Presence))
// }

// GetAllPresences retrieves all presence records
func GetAllPresences() ([]Presence, error) {
	o := orm.NewOrm()
	var presences []Presence
	_, err := o.QueryTable(new(Presence)).All(&presences)
	return presences, err
}

// GetPresenceById retrieves a presence record by ID
func GetPresenceById(id int) (*Presence, error) {
	o := orm.NewOrm()
	presence := Presence{Id: id}
	err := o.Read(&presence)
	if err == orm.ErrNoRows {
		return nil, nil
	}
	return &presence, err
}

// CreatePresence inserts a new presence record
func CreatePresence(p *Presence) error {
	o := orm.NewOrm()
	_, err := o.Insert(p)
	return err
}

// UpdatePresence updates an existing presence record
func UpdatePresence(p *Presence) error {
	o := orm.NewOrm()
	_, err := o.Update(p)
	return err
}

// DeletePresence deletes a presence record by ID
func DeletePresence(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&Presence{Id: id})
	return err
}
