package seeds

import (
	"log"

	"go-get-it/app/models"

	"gorm.io/gorm"
)

type SeedOption struct {
	ClearDB    bool `deafult:false`
	PopulateDb bool `deafult:false`
}

func down(db *gorm.DB) {
	err := db.Migrator().DropTable(&models.Post{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table")
	}
}

func Run(db *gorm.DB, option SeedOption) {
	if option.ClearDB && db.Migrator().HasTable(&models.User{}) {
		down(db)
	}

	if option.PopulateDb {

		for i, _ := range users {
			err := db.Model(&models.User{}).Create(&users[i]).Error
			if err != nil {
				log.Fatalf("cannot seed users table: %v", err)
			}
			posts[i].AuthorID = users[i].ID

			err = db.Model(&models.Post{}).Create(&posts[i]).Error
			if err != nil {
				log.Fatalf("cannot seed posts table: %v", err)
			}
		}
	}
}
