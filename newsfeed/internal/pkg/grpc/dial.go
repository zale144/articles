package grpc

import (
	c "articles/newsfeed/internal/config"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	v "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

func DialGrpc() (*grpc.ClientConn, error) {

	retriableErrors := []codes.Code{codes.DataLoss}
	retryTimeout := 50 * time.Millisecond

	unaryInterceptor := grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithCodes(retriableErrors...),
		grpc_retry.WithMax(3),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(retryTimeout)),
	)

	grpcConn, err := grpc.Dial(v.GetString(c.GrpcHost) + ":" + v.GetString(c.GrpcPort),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		return nil, err
	}

	log.Println("Dial gRPC at", v.GetString(c.GrpcPort))
	return grpcConn, nil
}
