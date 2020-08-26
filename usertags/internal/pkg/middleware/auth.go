package middleware

import (
	c "articles/usertags/internal/config"
	"articles/usertags/internal/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	v "github.com/spf13/viper"
	"time"
)

func GetUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*dto.JWTCustomClaims)
		c.Set("email", claims.Email)
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func GetJWTToken(email string) (string, error) {
	claims := &dto.JWTCustomClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(v.GetString(c.JWTSecret)))
	if err != nil {
		return "", err
	}
	return t, nil
}
