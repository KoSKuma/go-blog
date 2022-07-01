package repository

import (
	"errors"

	"github.com/KoSKuma/go-blog/api/adapter"
	"github.com/KoSKuma/go-blog/api/entity"
	"github.com/goccy/go-json"
)

type Logger interface {
	Log(route string, item interface{}) error
}

type Log struct {
	LogAdapter adapter.QueueAdapter
}

func (l *Log) Log(route string, item interface{}) error {
	message := ""
	switch item.(type) {
	case string:
		message = item.(string)
	case []entity.Post:
		m, _ := item.([]entity.Post)
		bm, err := json.Marshal(m)
		if err != nil {
			return err
		}
		message = string(bm)
	case entity.Post:
		m, _ := item.(entity.Post)
		bm, err := json.Marshal(m)
		if err != nil {
			return err
		}
		message = string(bm)
	case entity.PostUpdate:
		m, _ := item.(entity.PostUpdate)
		bm, err := json.Marshal(m)
		if err != nil {
			return err
		}
		message = string(bm)
	default:
		return errors.New("Invalid item type")
	}

	logItem := make(map[string]string)
	logItem["route"] = route
	logItem["message"] = message
	logMessage, err := json.Marshal(logItem)
	if err != nil {
		return err
	}
	return l.LogAdapter.Publish(string(logMessage))
}
