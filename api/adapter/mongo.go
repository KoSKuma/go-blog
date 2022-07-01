package adapter

import (
	"context"
	"errors"
	"time"

	"github.com/KoSKuma/go-blog/api/entity"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseAdapter interface {
	InsertOne(document entity.Post) (string, error)
	FindOne(id string) (entity.Post, error)
	UpdateOne(id string, document entity.PostUpdate) error
	DeleteOne(id string) error
	FindAll() ([]entity.Post, error)
}

type MongoAdapter struct {
	username   string
	password   string
	host       string
	port       string
	database   string
	collection string
}

func (m *MongoAdapter) Setup(username, password, host, port, database, collection string) {
	m.username = username
	m.password = password
	m.host = host
	m.port = port
	m.database = database
	m.collection = collection
}

func (m *MongoAdapter) InsertOne(document entity.Post) (string, error) {
	gId := uuid.New().String()
	document.Id = gId

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.username+":"+m.password+"@"+m.host+":"+m.port+"/admin"))
	if err != nil {
		return "", err
	}
	defer client.Disconnect(ctx)

	_, err = client.Database(m.database).Collection(m.collection).InsertOne(ctx, document)

	if err != nil {
		return "0", err
	}
	return gId, nil
}

func (m *MongoAdapter) FindOne(id string) (entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.username+":"+m.password+"@"+m.host+":"+m.port+"/admin"))
	if err != nil {
		return entity.Post{}, err
	}
	defer client.Disconnect(ctx)

	var post entity.Post
	err = client.Database(m.database).Collection(m.collection).FindOne(ctx, bson.M{"id": id}).Decode(&post)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return entity.Post{}, errors.New("Post not found")
		} else {
			return entity.Post{}, err
		}
	}
	return post, nil
}

func (m *MongoAdapter) UpdateOne(id string, document entity.PostUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.username+":"+m.password+"@"+m.host+":"+m.port+"/admin"))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	_, err = client.Database(m.database).Collection(m.collection).UpdateOne(ctx, bson.M{"id": id}, bson.D{
		primitive.E{Key: "$set", Value: document},
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoAdapter) DeleteOne(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.username+":"+m.password+"@"+m.host+":"+m.port+"/admin"))
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	_, err = client.Database(m.database).Collection(m.collection).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoAdapter) FindAll() ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+m.username+":"+m.password+"@"+m.host+":"+m.port+"/admin"))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	var posts []entity.Post
	result, err := client.Database(m.database).Collection(m.collection).Find(ctx, bson.M{})
	result.All(ctx, &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
