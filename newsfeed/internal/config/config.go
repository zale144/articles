package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const (
	configPath = "NF_CONFIG"
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
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("NF")
	viper.AutomaticEnv()
	path := "internal/config/config.yaml"

	if ep := os.Getenv(configPath); len(ep) > 0 {
		path = ep
	}

	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		log.Println(err)
	}

	return viper.MergeInConfig()
}
