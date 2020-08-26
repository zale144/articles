package main

import (
	"articles/newsfeed/internal/api/handler"
	"articles/newsfeed/internal/api/server"
	"articles/newsfeed/internal/config"
	"articles/newsfeed/internal/pkg/db"
	g "articles/newsfeed/internal/pkg/grpc"
	"articles/newsfeed/internal/repository"
	cService "articles/newsfeed/internal/service/cards"
	"context"
	"github.com/zale144/articles/pb"
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
		if err := dbc.Disconnect(context.TODO()); err != nil {
			log.Println("error while closing the database connection", err)
		}
	}()

	grpcConn, err := g.DialGrpc()
	if err != nil {
		log.Fatal(err)
	}

	store := repository.NewStore(dbc)
	client := pb.NewTagsServiceClient(grpcConn)

	cSvc := cService.NewCardsService(store, client)

	server.Run(handler.NewHandler(cSvc))
}
