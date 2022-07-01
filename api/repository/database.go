package repository

import (
	"github.com/KoSKuma/go-blog/api/adapter"
	"github.com/KoSKuma/go-blog/api/entity"
)

type DatabaseRepo interface {
	InsertOne(document entity.Post) (string, error)
	FindOne(id string) (entity.Post, error)
	UpdateOne(id string, document entity.PostUpdate) error
	DeleteOne(id string) error
	FindAll() ([]entity.Post, error)
}

type Database struct {
	DBAdapter adapter.DatabaseAdapter
}

func (d *Database) InsertOne(document entity.Post) (string, error) {
	return d.DBAdapter.InsertOne(document)
}

func (d *Database) FindOne(id string) (entity.Post, error) {
	return d.DBAdapter.FindOne(id)
}

func (d *Database) UpdateOne(id string, document entity.PostUpdate) error {
	return d.DBAdapter.UpdateOne(id, document)
}

func (d *Database) DeleteOne(id string) error {
	return d.DBAdapter.DeleteOne(id)
}

func (d *Database) FindAll() ([]entity.Post, error) {
	return d.DBAdapter.FindAll()
}
