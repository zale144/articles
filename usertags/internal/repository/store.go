package repository

import (
	"articles/usertags/internal/model"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return Store{
		db: db,
	}
}

func (u Store) AddTagsToUser(user *model.User) error {

	lenT := len(user.Tags)
	if lenT == 0 {
		return errors.New("no tags provided")
	}

	var valueArgs []interface{}
	for _, tag := range user.Tags {
		valueArgs = append(valueArgs, tag.Keyword, "NOW()", "NOW()")
	}
	valueArgs = append(valueArgs, user.ID)

	colsT := []string{"keyword", "created_at", "updated_at"}
	colsUT := []string{"tag_id, user_id"}

	stmt := fmt.Sprintf(
		`WITH rows AS ( %s ON CONFLICT (keyword) DO UPDATE SET updated_at = NOW() RETURNING id) %s SELECT rows.id, $%d FROM rows`,
		buildInsertStr(lenT, "tag", colsT), buildInsertStr(0, "user_tags", colsUT), lenT*len(colsT)+1)

	_, err := u.db.DB().Exec(stmt, valueArgs...)
	return err
}

func (u Store) GetUser(email string, withTags bool) (*model.User, error) {
	user := &model.User{}

	if withTags {
		if err := u.db.Preload("Tags").Where("email = ?", email).First(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}

	return user, u.db.Where("email = ?", email).First(user).Error
}

func (u Store) CreateUser(user *model.User) error {
	return u.db.Create(user).Error
}
