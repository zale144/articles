package db

import (
	c "articles/usertags/internal/config"
	"articles/usertags/internal/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	v "github.com/spf13/viper"
)

func Setup() (*gorm.DB, error) {

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		v.GetString(c.DBHost), v.GetString(c.DBPort), v.GetString(c.DBName), v.GetString(c.DBUser), v.GetString(c.DBPassword)))
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.User{},
		&model.Tag{},
	)

	db.Model(&model.UserTags{}).
		AddForeignKey(`"user_id"`, `"user"(id)`, "RESTRICT", "RESTRICT").
		AddForeignKey(`"tag_id"`, `"tag"(id)`, "RESTRICT", "RESTRICT")

	return db, nil
}
