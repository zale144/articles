package model

import "github.com/jinzhu/gorm"

type Tag struct {
	gorm.Model
	Keyword   string `gorm:"keyword:;unique_index"`

}

// TableName returns name of table tag
func (u Tag) TableName() string {
	return "tag"
}
