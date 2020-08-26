package middleware

import (
	"articles/newsfeed/internal/dto"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
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
