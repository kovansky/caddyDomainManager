package databases

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoSource struct {
	User     string
	Password string
	Host     string
	Port     int
	AuthDb   string

	client        *mongo.Client
	database      *mongo.Database
	connectionUri string
}

func (source *MongoSource) Connect() bool {
	source.BuildUri()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(source.connectionUri))
	if err != nil {
		return false
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false
	}

	source.client = client

	return true
}

func (source *MongoSource) BuildUri() {
	uri := "mongodb://"

	if len(source.User) > 0 {
		uri += source.User

		if len(source.Password) > 0 {
			uri += fmt.Sprintf(":%s", source.Password)
		}

		uri += "@"
	}

	if len(source.Host) > 0 {
		uri += source.Host

		if source.Port != 0 {
			uri += fmt.Sprintf(":%d", source.Port)
		}
	}

	if len(source.AuthDb) > 0 {
		uri += fmt.Sprintf("/?authDatabase=%s", source.AuthDb)
	}

	source.connectionUri = uri
}

func (source MongoSource) CreateUser(name string, userHost string, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := source.database.RunCommand(ctx, bson.D{
		{"createUser", name},
		{"pwd", password},
		{"roles", []bson.M{{"role": "readWrite", "db": source.database.Name()}}},
		{"authenticationRestrictions", []bson.M{{"clientSource": bson.A{userHost}}}},
	})

	if result.Err() != nil {
		return false
	}

	return true
}

func (source *MongoSource) CreateDatabase(name string) bool {
	// Mongo does not need database creation

	source.database = source.client.Database(name)

	return true
}

func (source MongoSource) Close() {
	_ = source.client.Disconnect(context.Background())
}
