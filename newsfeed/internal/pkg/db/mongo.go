package db

import (
	c "articles/newsfeed/internal/config"
	"context"
	"fmt"
	v "github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

func Setup() (*mongo.Client, error) {

	var cred options.Credential

	cred.Username = v.GetString(c.DBUser)
	cred.Password = v.GetString(c.DBPassword)

	connStr := fmt.Sprintf("mongodb://%s:%s", v.GetString(c.DBHost), v.GetString(c.DBPort))
	clientOptions := options.Client().ApplyURI(connStr).SetAuth(cred)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
