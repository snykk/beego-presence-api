package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id         int         `orm:"auto" json:"id"`
	Name       string      `orm:"size(100)" json:"name"`
	Email      string      `orm:"size(100);unique" json:"email"`
	Password   string      `orm:"size(255)" json:"-"`
	Role       string      `orm:"size(10)" json:"role"`
	Department *Department `orm:"rel(fk);column(department_id)" json:"department"`  // ForeignKey to Department
	Presences  []*Presence `orm:"reverse(many)" json:"presences"`                   // Reverse relationship with Presence
	Schedule   *Schedule   `orm:"null;rel(fk);column(schedule_id)" json:"schedule"` // ForeignKey to Schedule
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)" json:"updated_at"`
}

// func init() {
// 	orm.RegisterModel(new(User))
// }

func GetAllUsers(isIncludePresenceList bool) ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	_, err := o.QueryTable(new(User)).RelatedSel("Department", "Schedule").All(&users)
	if err != nil {
		return nil, err
	}

	if isIncludePresenceList {
		for i := range users {
			_, err := o.LoadRelated(users[i], "Presences")
			if err != nil {
				return nil, err
			}
		}
	}

	return users, nil
}

func GetUserByEmail(email string) (User, error) {
	o := orm.NewOrm()
	user := User{Email: email}
	err := o.Read(&user, "Email")
	return user, err
}

func GetUserById(id int, isIncludePresenceList bool) (*User, error) {
	o := orm.NewOrm()
	user := &User{Id: id}
	err := o.Read(user)
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(user, "Department")
	if err != nil {
		return nil, err
	}

	_, err = o.LoadRelated(user, "Schedule")
	if err != nil {
		return nil, err
	}

	if isIncludePresenceList {
		_, err := o.LoadRelated(user, "Presences")
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func CreateUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(user)
	return err
}

func UpdateUser(user *User) error {
	o := orm.NewOrm()
	_, err := o.Update(user)
	return err
}

func DeleteUser(id int) error {
	o := orm.NewOrm()
	_, err := o.Delete(&User{Id: id})
	return err
}
