package main

import (
	"github.com/KoSKuma/go-blog/log_worker/adapter"
	"github.com/KoSKuma/go-blog/log_worker/usecase"
)

func main() {

	mongoDatabaseAdapter := adapter.MongoAdapter{Username: "root", Password: "root", Host: "localhost", Port: "27017", Database: "blog", Collection: "logs"}
	rabbitMQAdapter := adapter.RabbitMQAdapter{Username: "root", Password: "root", Host: "localhost", Port: "5672", Queue: "blog:logs"}
	bulkUsecase := usecase.BulkLogUsercase{Bulker: &mongoDatabaseAdapter, Subscriber: &rabbitMQAdapter}

	bulkUsecase.BulkLog()
}
