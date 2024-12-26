package models

import (
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type User struct {
	Id         int         `orm:"auto"`
	Name       string      `orm:"size(100)"`
	Email      string      `orm:"size(100);unique"`
	Password   string      `orm:"size(255)"`
	Department *Department `orm:"rel(fk);column(department_id)"` // ForeignKey to Department
	Presences  []*Presence `orm:"reverse(many)"`                 // Reverse relationship with Presence
	CreatedAt  time.Time   `orm:"auto_now_add;type(datetime)"`
	UpdatedAt  time.Time   `orm:"auto_now;type(datetime)"`
}

// func init() {
// 	orm.RegisterModel(new(User))
// }

func GetAllUsers() ([]User, error) {
	o := orm.NewOrm()
	var users []User
	_, err := o.QueryTable(new(User)).All(&users)
	return users, err
}

func GetUserById(id int) (User, error) {
	o := orm.NewOrm()
	user := User{Id: id}
	err := o.Read(&user)
	return user, err
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
