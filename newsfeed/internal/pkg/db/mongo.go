package db

import (
	c "articles/newsfeed/internal/config"
	"context"
	"fmt"
	v "github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

// Setup initiates DB connection. If success it will return a pointer to DB instance.
func Setup() (*mongo.Client, error) {

	var cred options.Credential

	cred.Username = v.GetString(c.DBUser)
	cred.Password = v.GetString(c.DBPassword)

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", v.GetString(c.DBHost), v.GetString(c.DBPort))).SetAuth(cred)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
