package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// Presence represents the presence table in the database
type Presence struct {
	Id        int       `orm:"auto" json:"id"`
	User      *User     `orm:"rel(fk)" json:"user"`     // ForeignKey to User
	Schedule  *Schedule `orm:"rel(fk)" json:"schedule"` // ForeignKey to Schedule
	Type      string    `orm:"size(10)" json:"type"`
	Status    string    `orm:"size(50)" json:"status"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// func init() {
// 	orm.RegisterModel(new(Presence))
// }

// GetAllPresences retrieves all presence records
func GetAllPresences() ([]*Presence, error) {
	o := orm.NewOrm()
	var presences []*Presence
	_, err := o.QueryTable(new(Presence)).RelatedSel("User", "Schedule").All(&presences)
	return presences, err
}

// GetPresenceById retrieves a presence record by ID
func GetPresenceById(id int) (*Presence, error) {
	o := orm.NewOrm()
	presence := &Presence{Id: id}
	err := o.Read(presence)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(presence, "User")
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(presence, "Schedule")
	if err != nil {
		return nil, err
	}

	return presence, err
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
