package dto

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

type GetCardsRequest struct {
	Tags []string `json:"tags" form:"tags" query:"tags" validate:"eq=2"`
}

type GetCardsPayload struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	Title     string    `json:"title" validate:"required"`
	Timestamp time.Time `json:"timestamp" validate:"required"`
	Tags      []string  `json:"tags" validate:"required"`
}

type JWTCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
