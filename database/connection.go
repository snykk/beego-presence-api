package database

import (
	"fmt"

	"github.com/snykk/beego-presence-api/models"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/server/web"
	_ "github.com/lib/pq"
)

func InitDB() {
	// Load configuration
	user := web.AppConfig.DefaultString("pg_user", "postgres")
	password := web.AppConfig.DefaultString("pg_password", "password")
	dbname := web.AppConfig.DefaultString("pg_dbname", "presence_db")
	host := web.AppConfig.DefaultString("pg_host", "localhost")
	port := web.AppConfig.DefaultString("pg_port", "5432")

	// Connection string
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	// Register driver and database
	orm.RegisterDriver("postgres", orm.DRPostgres)
	err := orm.RegisterDataBase("default", "postgres", dsn)
	if err != nil {
		panic(err)
	}

	// Register Models
	orm.RegisterModel(new(models.User), new(models.Department), new(models.Schedule), new(models.Presence))

	// Auto Create Tables
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("Daatbase connected successfully!")
}
