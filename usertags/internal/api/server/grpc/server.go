package grpc

import (
	"articles/usertags/internal/api/handler/tags"
	c "articles/usertags/internal/config"
	"github.com/pkg/errors"
	v "github.com/spf13/viper"
	"github.com/zale144/articles/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Start(tagsHnd tags.Tags) error {

	url := v.GetString(c.GrpcHost) + ":" + v.GetString(c.GrpcPort)
	lis, err := net.Listen("tcp", url)
	if err != nil {
		return errors.Wrap(err, "failed to listen gRPC")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTagsServiceServer(grpcServer, tagsHnd)

	log.Printf("gRPC listening on: %s", url)
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve gRPC")
	}
	return nil
}
