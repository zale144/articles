package http

import (
	"articles/usertags/internal/api/handler"
	c "articles/usertags/internal/config"
	"articles/usertags/internal/dto"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v "github.com/spf13/viper"
	"log"

	mware "articles/usertags/internal/pkg/middleware"
)

func Run(hnd handler.Handler)  {

	e := echo.New()
	e.Debug = true
	e.Use(middleware.CORS())
	e.Use(middleware.Secure())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	valid := validator.New()
	e.Validator = mware.NewCustomValidator(valid)

	api := e.Group("/api")

	a := api.Group("/a")
	a.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &dto.JWTCustomClaims{},
		SigningKey: []byte(v.GetString(c.JWTSecret)),
	}))

	a.Use(mware.GetUser)

	api.POST("/register", hnd.Register)
	api.POST("/login", hnd.Login)

	t := a.Group("/tag")

	t.POST("/", hnd.AddTag)
	t.GET("/", hnd.GetTag)

	log.Fatal(e.Start(":" + v.GetString(c.HttpPort)))
}
