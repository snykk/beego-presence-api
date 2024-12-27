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

// GetPresencesByUserId retrieves all presence records for a given user ID
func GetPresencesByUserId(userId int) ([]*Presence, error) {
	o := orm.NewOrm()
	var presences []*Presence
	_, err := o.QueryTable(new(Presence)).
		Filter("User__Id", userId).
		RelatedSel("User", "Schedule").
		All(&presences)
	return presences, err
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
func DeletePresence(id int) (int64, error) {
	o := orm.NewOrm()
	affectedRows, err := o.Delete(&Presence{Id: id})
	return affectedRows, err
}

// CheckPresenceExistsByUserAndType checks if a presence record exists for a given user ID, presence type, and date
func CheckPresenceExistsByUserAndType(userId int, presenceType string, date time.Time) (bool, error) {
	// Ensure the date is in the same timezone as the database
	location, _ := time.LoadLocation("Asia/Jakarta")
	date = date.In(location)

	o := orm.NewOrm()
	count, err := o.QueryTable(new(Presence)).
		Filter("User__Id", userId).
		Filter("Type", presenceType).
		Filter("CreatedAt__gte", date.Format("2006-01-02")+" 00:00:00").
		Filter("CreatedAt__lte", date.Format("2006-01-02")+" 23:59:59").
		Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
