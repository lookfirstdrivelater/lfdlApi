package lfdlapi

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	// there is a blank import
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Event struct {
	gorm.Model
	StartTime time.Time
	EndTime   time.Time
	Points  string
	EventType string
	Severity  string
	CenterX string
	CenterY string
}

type user struct {
	gorm.Model
	UserName    string
	FirstName   string
	LastName    string
	Password    string
	AuthGeneral bool
	AuthAdmin   bool
}

var db *gorm.DB

// API is the main entry point for this lib
func API(DB *gorm.DB) {

	db = DB

	// run migrations only if the tables do not exist
	//if !DB.HasTable(event{}) {
		DB.AutoMigrate(Event{})
		log.Println("Migrating Event")
	//}
	if !DB.HasTable(user{}) {
		DB.AutoMigrate(user{})
		log.Println("Migrating Users")

		pass, err := hashPassword(os.Getenv("DatabaseDefaultAdminPassword"))
		isGeneral, _ := strconv.ParseBool(os.Getenv("DatabaseDefaultAdminLevelGeneral"))
		isAdmin, _  := strconv.ParseBool(os.Getenv("DatabaseDefaultAdminLevelAdmin"))

		fmt.Println(err)

		db.Create(&user{
			UserName:  os.Getenv("DatabaseDefaultAdminUserName"),
			FirstName: os.Getenv("DatabaseDefaultAdminFirstName"),
			LastName:  os.Getenv("DatabaseDefaultAdminLastName"),
			Password:  pass,
			AuthGeneral: isGeneral,
			AuthAdmin: isAdmin,
		})
	}

	// create restful router
	router := setupRouter()
	// start the server
	_ = router.Run(fmt.Sprintf(":8080"))
}
