package config

import (
	"github.com/spf13/viper"
	"log"
)

const (
	JWTSecret = "jwtsecret"
	HttpPort  = "http.port"
	GrpcHost  = "grpc.host"
	GrpcPort  = "grpc.port"

	DBHost     = "db.host"
	DBPort     = "db.port"
	DBUser     = "db.user"
	DBPassword = "db.password"
	DBName     = "db.name"
)

func Configure() error {
	viper.SetConfigFile("internal/config/config.yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		return err
	}

	return nil
}
