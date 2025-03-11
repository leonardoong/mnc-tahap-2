package queue

import (
	"github.com/leonardoong/e-wallet/config"
	"github.com/leonardoong/e-wallet/internal/consumer"
	"github.com/leonardoong/e-wallet/internal/service"
)

type Queue struct {
	Consumer *consumer.Consumer
}

func NewQueue(cfg *config.Config, svc service.ITransactionService) *Queue {
	queue := new(Queue)
	queue.Consumer = consumer.NewConsumer(cfg, svc)
	return queue
}

func (q *Queue) Initialize() {
	q.Consumer.Initialize()
}
