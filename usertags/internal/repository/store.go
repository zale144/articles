package repository

import (
	"articles/usertags/internal/model"
	"articles/usertags/internal/pkg/db"
	"database/sql"
	"errors"
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

	if len(user.Tags) == 0 {
		return errors.New("no tags provided")
	}

	return db.Transact(u.db.DB(), func(tx *sql.Tx) (e error) {

		now := time.Now()

		var valueArgs []interface{}
		for _, tag := range user.Tags {
			valueArgs = append(valueArgs, tag.Keyword, now, now)
		}

		cols := []string{"keyword", "created_at", "updated_at"}
		stmt := fmt.Sprintf("%s ON CONFLICT (keyword) DO UPDATE SET keyword = EXCLUDED.keyword, updated_at = NOW() RETURNING id",
			prepareBulkInsertStmt(len(user.Tags), "tag", cols))

		rows, e := tx.Query(stmt, valueArgs...)
		if e != nil {
			return
		}

		var ids []int
		for rows.Next() {
			id := 0
			if e = rows.Scan(&id); e != nil {
				return
			}
			ids = append(ids, id)
		}

		valueArgs = []interface{}{}

		for _, tagID := range ids {
			valueArgs = append(valueArgs, user.ID, tagID)
		}

		cols = []string{"user_id", "tag_id"}
		stmt = fmt.Sprintf("%s ON CONFLICT DO NOTHING", prepareBulkInsertStmt(len(ids), "user_tags", cols))

		if _, e = tx.Exec(stmt, valueArgs...); e != nil {
			return
		}
		return
	})
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

func prepareBulkInsertStmt(rowsL int, tableName string, cols []string) string {

	valueStrings := make([]string, 0, rowsL)
	valueArgs := make([]interface{}, 0, rowsL*len(cols))
	colsN := len(cols)

	for i := 0; i < rowsL; i++ {

		var placeholders []string

		for j, col := range cols {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*colsN+j+1))
			valueArgs = append(valueArgs, col)
		}

		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(placeholders, ",")))
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", tableName, strings.Join(cols, ","), strings.Join(valueStrings, ","))
}
