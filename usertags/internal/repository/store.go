package repository

import (
	"articles/usertags/internal/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
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

	tx := u.db.Begin()

	valueStrings := []string{}
	valueArgs := []interface{}{}
	now := time.Now()

	for _, tag := range user.Tags {
		valueStrings = append(valueStrings, "(?,?,?)")
		valueArgs = append(valueArgs, tag.Keyword)
		valueArgs = append(valueArgs, now)
		valueArgs = append(valueArgs, now)
	}

	stmt := fmt.Sprintf("INSERT INTO tag (keyword, created_at, updated_at) VALUES %s ON CONFLICT DO NOTHING", strings.Join(valueStrings, ","))

	if err := tx.Exec(stmt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}

	words := make([]string, len(user.Tags))

	for i := range user.Tags {
		words[i] = user.Tags[i].Keyword
	}

	var tags []model.Tag
	if err := tx.Where("keyword IN (?)", words).Find(&tags).Error; err != nil {
		tx.Rollback()
		return err
	}

	valueStrings = []string{}
	valueArgs = []interface{}{}

	for _, tag := range tags {
		valueStrings = append(valueStrings, "(?,?)")
		valueArgs = append(valueArgs, user.ID)
		valueArgs = append(valueArgs, tag.ID)
	}

	stmt = fmt.Sprintf("INSERT INTO user_tags (user_id, tag_id) VALUES %s ON CONFLICT DO NOTHING", strings.Join(valueStrings, ","))

	if err := tx.Exec(stmt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (u Store) GetUser(email string, withTags bool) (*model.User, error) {
	user := &model.User{}

	if withTags {
		if err := u.db.Preload("Tags").Where("email = ?", email).First(user).Error; err != nil {
			return nil, err
		}
		return user, nil
	}

	if err := u.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u Store) CreateUser(user *model.User) error {
	return u.db.Create(user).Error
}
