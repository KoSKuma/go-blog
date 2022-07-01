package usecase

import (
	"fmt"
	"time"

	"github.com/KoSKuma/go-blog/log_worker/adapter"
	amqp "github.com/rabbitmq/amqp091-go"
)

type BulkLogUsercase struct {
	Bulker     adapter.Bulker
	Subscriber adapter.Subscriber
}

func (b *BulkLogUsercase) BulkLog() {
	tick := time.Tick(1 * time.Minute)
	b.Subscriber.Subscribe(b.Callback, tick)
}

func (b *BulkLogUsercase) Callback(msgs <-chan amqp.Delivery, tick <-chan time.Time) error {
	logs := []amqp.Delivery{}
	for {
		select {
		case d := <-msgs:
			logs = append(logs, d)
			if len(logs) >= 10 {
				fmt.Println(len(logs))
				formattedLogs := formatLog(logs)
				b.Bulker.BulkInsert(formattedLogs)
				d.Ack(true)
				logs = []amqp.Delivery{}
			}
		case <-tick:
			if len(logs) > 0 {
				fmt.Println(len(logs))
				formattedLogs := formatLog(logs)
				b.Bulker.BulkInsert(formattedLogs)
				logs[len(logs)-1].Ack(true)
				logs = []amqp.Delivery{}
			}
		}
	}
}

func formatLog(dLogs []amqp.Delivery) []string {
	logs := []string{}
	for _, d := range dLogs {
		logs = append(logs, string(d.Body))
	}
	return logs
}
