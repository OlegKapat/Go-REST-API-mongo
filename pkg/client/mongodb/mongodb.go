package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDb string) (db *mongo.Database, err error) {
	var isAuth bool
	var Mongo_URI string
	if username == "" && password == "" {
		Mongo_URI = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		Mongo_URI = fmt.Sprintf("mongodb://%s:%s@%s:$s", username, password, host, port)
	}

	clientOptions := options.Client().ApplyURI(Mongo_URI)
	if isAuth {
		if authDb == "" {
			authDb = database
		}
		clientOptions.SetAuth(options.Credential{
			AuthSource: authDb,
			Username:   username,
			Password:   password,
		})
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to mongoDb due to error %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Failed to ping error %/v", err)
	}
	return client.Database(database), nil
}
