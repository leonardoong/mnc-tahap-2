package consumer

import (
	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/service"
)

type topUpWorker struct {
	transactionService service.ITransactionService
	workerPool         *work.WorkerPool
	jobName            string
}

func newTopUpWorker(srv service.ITransactionService, pool *work.WorkerPool) *topUpWorker {
	return &topUpWorker{
		transactionService: srv,
		workerPool:         pool,
	}
}

func (c *topUpWorker) runTopupConsumer(maxFails uint) {
	c.workerPool.JobWithOptions(c.jobName, work.JobOptions{MaxFails: maxFails}, c.processTopUp)
}

func (c *topUpWorker) processTopUp(job *work.Job) (err error) {
	if err = job.ArgError(); err != nil {
		return err
	}

	req := entity.PublishTopUpRequest{
		Amount:  job.ArgFloat64("amount"),
		TopUpID: job.ArgString("top_up_id"),
		UserID:  job.ArgString("user_id"),
	}
	err = c.transactionService.ProcessTopUp(req)
	if err != nil {
		return
	}
	return
}
