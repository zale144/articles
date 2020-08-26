package service

import (
	"articles/usertags/internal/dto"
	"articles/usertags/internal/model"
	"articles/usertags/internal/pkg/middleware"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	store store
}

func NewUserService(store store) User {
	return User{
		store: store,
	}
}

type store interface {
	CreateUser(*model.User) error
	GetUser(string, bool) (*model.User, error)
}

func (r User) Register(userDto dto.RegisterPayload) error {

	password, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:     userDto.Name,
		Email:    userDto.Email,
		Password: string(password),
	}

	return r.store.CreateUser(user)
}

func (r User) Login(userDto dto.LoginPayload) (string, error) {

	user, err := r.store.GetUser(userDto.Email, true)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password)); err != nil {
		return "", err
	}

	return middleware.GetJWTToken(userDto.Email)
}
