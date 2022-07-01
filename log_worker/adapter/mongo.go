package adapter

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Bulker interface {
	BulkInsert(documents []string) error
}

type MongoAdapter struct {
	Username   string
	Password   string
	Host       string
	Port       string
	Database   string
	Collection string
}

func (m *MongoAdapter) BulkInsert(documents []string) error {
	models := []mongo.WriteModel{}
	for _, document := range documents {
		models = append(models, mongo.NewInsertOneModel().SetDocument(bson.M{"log": document}))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.Username+":"+m.Password+"@"+m.Host+":"+m.Port+"/admin"))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	_, err = client.Database(m.Database).Collection(m.Collection).BulkWrite(ctx, models, nil)
	if err != nil {
		return err
	}
	return nil
}
