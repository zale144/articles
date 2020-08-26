package server

import (
	"articles/newsfeed/internal/api/handler"
	c "articles/newsfeed/internal/config"
	"articles/newsfeed/internal/dto"
	mware "articles/newsfeed/internal/pkg/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v "github.com/spf13/viper"
	"log"
)

func Run(hnd handler.Handler) {

	e := echo.New()
	e.Debug = true // TODO
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

	t := a.Group("/card")

	t.POST("/", hnd.AddCard)
	t.GET("/", hnd.GetCards)
	t.GET("/by-tags", hnd.GetCardsByTags)

	log.Fatal(e.Start(":" + v.GetString(c.HttpPort)))
}
