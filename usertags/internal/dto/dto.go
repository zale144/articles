package dto

import "github.com/dgrijalva/jwt-go"

type RegisterPayload struct {
	Name string `json:"name" validate:"required"`
	LoginPayload
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ResponseMessage struct {
	Message string `json:"message"`
}

type AddTagsPayload struct {
	Tags []string `json:"tags" validate:"required"`
}

type GetTagsPayload struct {
	Tags []string `json:"tags"`
}

type JWTCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
