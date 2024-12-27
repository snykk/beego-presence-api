package database

import (
	"fmt"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/snykk/beego-presence-api/constants"
	"github.com/snykk/beego-presence-api/helpers"
	"github.com/snykk/beego-presence-api/models"
)

// SeedDepartments populates the departments table if it's empty
func SeedDepartments() {
	o := orm.NewOrm()

	// Check if the departments table is empty
	count, err := o.QueryTable(new(models.Department)).Count()
	if err != nil {
		log.Printf("Failed to check departments table: %v", err)
		return
	}

	if count == 0 {
		departments := []models.Department{
			{Name: "Engineering"},
			{Name: "Human Resources"},
			{Name: "Marketing"},
		}

		for _, dept := range departments {
			_, err := o.Insert(&dept)
			if err != nil {
				log.Printf("Failed to seed department %s: %v", dept.Name, err)
			} else {
				fmt.Printf("Seeded department: %s\n", dept.Name)
			}
		}
	} else {
		fmt.Println("Departments table already seeded.")
	}
}

// SeedUsers populates the users table if it's empty
func SeedUsers() {
	o := orm.NewOrm()

	// Check if the users table is empty
	count, err := o.QueryTable(new(models.User)).Count()
	if err != nil {
		log.Printf("Failed to check users table: %v", err)
		return
	}

	if count == 0 {
		// Ensure departments exist
		departmentNames := []string{"Engineering", "Human Resources", "Marketing"}
		var departments []models.Department

		for _, deptName := range departmentNames {
			dept := models.Department{Name: deptName}
			if created, _, err := o.ReadOrCreate(&dept, "Name"); err != nil {
				log.Printf("Failed to ensure department %s exists: %v", deptName, err)
			} else {
				if created {
					fmt.Printf("Created department: %s\n", dept.Name)
				} else {
					fmt.Printf("Department already exists: %s\n", dept.Name)
				}
			}
			departments = append(departments, dept)
		}

		// Prepare users
		users := []models.User{
			{Name: "Admin", Role: constants.RoleAdmin, Email: "admin@example.com", Password: "1234", Department: &departments[0]},
			{Name: "Employee1", Role: constants.RoleEmployee, Email: "employee1@example.com", Password: "1234", Department: &departments[1]},
			{Name: "Employee2", Role: constants.RoleEmployee, Email: "employee2@example.com", Password: "1234", Department: &departments[2]},
		}

		// Seed users and assign schedules
		for _, user := range users {
			hashedPassword, err := helpers.HashPassword(user.Password)
			if err != nil {
				log.Fatalf("Error hashing password %v\n", err)
			}
			user.Password = hashedPassword

			if user.Role == constants.RoleEmployee {
				// Assign a schedule before inserting user
				err = assignScheduleToUser(&user, o)
				if err != nil {
					log.Printf("Failed to assign schedule to user %s: %v", user.Name, err)
					continue
				}
			}

			_, err = o.Insert(&user)
			if err != nil {
				log.Printf("Failed to seed user %s: %v", user.Name, err)
			} else {
				fmt.Printf("Seeded user: %s in department: %s\n", user.Name, user.Department.Name)
			}
		}
	} else {
		fmt.Println("Users table already seeded.")
	}
}

func assignScheduleToUser(user *models.User, o orm.Ormer) error {
	if user.Department == nil {
		return fmt.Errorf("user %s does not have a department", user.Name)
	}

	var schedules []models.Schedule
	_, err := o.QueryTable(new(models.Schedule)).Filter("Department", user.Department.Id).All(&schedules)
	if err != nil {
		return fmt.Errorf("failed to fetch schedules for department %d: %v", user.Department.Id, err)
	}

	if len(schedules) == 0 {
		return fmt.Errorf("no schedules available for department %d", user.Department.Id)
	}

	// Randomly select a schedule
	randomSchedule := schedules[helpers.RandomInt(0, len(schedules)-1)]
	user.Schedule = &randomSchedule
	return nil
}

func SeedSchedules() {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(models.Schedule)).Count()
	if err != nil {
		log.Printf("Failed to check schedules table: %v", err)
		return
	}

	if count == 0 {
		departmentNames := []string{"Engineering", "Human Resources", "Marketing"}
		var departments []models.Department

		for _, deptName := range departmentNames {
			dept := models.Department{Name: deptName}
			if err := o.Read(&dept, "Name"); err == nil {
				departments = append(departments, dept)
			} else {
				log.Printf("Failed to find department %s: %v", deptName, err)
			}
		}

		schedules := []models.Schedule{
			// Engineering Department
			{Name: "Software Development Shift", Department: &departments[0], InTime: "09:00:00", OutTime: "17:00:00"},
			{Name: "Maintenance and Support", Department: &departments[0], InTime: "18:00:00", OutTime: "02:00:00"},

			// Human Resources Department
			{Name: "Recruitment and Interviews", Department: &departments[1], InTime: "10:00:00", OutTime: "16:00:00"},
			{Name: "HR Administration", Department: &departments[1], InTime: "08:00:00", OutTime: "15:00:00"},

			// Marketing Department
			{Name: "Campaign Planning", Department: &departments[2], InTime: "09:30:00", OutTime: "17:30:00"},
			{Name: "Client Meetings and Networking", Department: &departments[2], InTime: "11:00:00", OutTime: "19:00:00"},
		}

		for _, schedule := range schedules {
			_, err := o.Insert(&schedule)
			if err != nil {
				log.Printf("Failed to seed schedule %s: %v", schedule.Name, err)
			} else {
				fmt.Printf("Seeded schedule: %s for department: %s\n", schedule.Name, schedule.Department.Name)
			}
		}
	} else {
		fmt.Println("Schedules table already seeded.")
	}
}

func RunAllSeeds() {
	SeedDepartments()
	SeedSchedules()
	SeedUsers()
}
