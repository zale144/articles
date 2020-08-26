package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email;unique_index"`
	Password string `gorm:"column:password"`
	Tags     []Tag  `gorm:"many2many:user_tags"`
}

// TableName returns name of table user
func (u User) TableName() string {
	return "user"
}
