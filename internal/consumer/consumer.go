package consumer

import (
	"fmt"
	"log"

	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/config"
	"github.com/leonardoong/e-wallet/internal/service"
)

type Consumer struct {
	config      *config.Config
	workerPool  *work.WorkerPool
	TopUpWorker *topUpWorker
	// PaymentWorker
	// TransferWorker
}

type WorkerContext struct{}

func NewConsumer(cfg *config.Config, svc service.ITransactionService) *Consumer {
	if svc == nil {
		log.Fatal("service.ITransactionService is nil in NewConsumer")
	}
	consumer := new(Consumer)
	consumer.config = cfg
	consumer.workerPool = work.NewWorkerPool(WorkerContext{}, uint(2), "ewallet", cfg.CachePool)
	consumer.TopUpWorker = newTopUpWorker(svc, consumer.workerPool)

	return consumer
}

func (c *Consumer) Initialize() {
	fmt.Println("INIT CONSUMER")
	maxFails := uint(2)

	c.TopUpWorker.workerPool = c.workerPool
	c.TopUpWorker.jobName = "topup_job"
	c.TopUpWorker.runTopupConsumer(maxFails)

	c.workerPool.Start()
}

func (c *Consumer) Close() {
	c.workerPool.Stop()
}
