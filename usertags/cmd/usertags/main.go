package main

import (
	"articles/usertags/internal/api/handler"
	"articles/usertags/internal/api/handler/tags"
	"articles/usertags/internal/api/server/grpc"
	"articles/usertags/internal/api/server/http"
	"articles/usertags/internal/config"
	"articles/usertags/internal/repository"
	tService "articles/usertags/internal/service/tags"
	uService "articles/usertags/internal/service/users"

	"articles/usertags/internal/pkg/db"
	"log"
)

func main() {

	if err := config.Configure(); err != nil {
		log.Fatal(err)
	}

	dbc, err := db.Setup()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbc.Close(); err != nil {
			log.Println("error while closing the database connection", err)
		}
	}()

	store := repository.NewStore(dbc)

	uSvc := uService.NewUserService(store)
	tSvc := tService.NewTagService(store)

	usrHnd := tags.NewTags(tSvc)

	go func() {
		if err := grpc.Start(usrHnd); err != nil {
			log.Fatal(err)
		}
	}()

	http.Run(handler.NewHandler(uSvc, tSvc))
}
