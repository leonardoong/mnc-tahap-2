package consumer

import (
	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/service"
)


type paymentWorker struct {
	transactionService service.ITransactionService
	workerPool *work.WorkerPool
	jobName string
}

func newPaymentWorker(srv service.ITransactionService, pool *work.WorkerPool) *paymentWorker {
	return &paymentWorker{
		transactionService: srv,
		workerPool:         pool,
	}
}

func (c *paymentWorker) runPaymentConsumer(maxFails uint) {
	c.workerPool.JobWithOptions(c.jobName, work.JobOptions{MaxFails: maxFails}, c.processPayment)
}

func (c *paymentWorker) processPayment(job *work.Job) (err error) {
	if err = job.ArgError(); err != nil {
		return err
	}

	req := entity.PaymentRequest{
		Amount:  job.ArgFloat64("amount"),
		PaymentID: job.ArgString("payment_id"),
		UserID:  job.ArgString("user_id"),
		Remarks: job.ArgString("remarks"),
	}
	err = c.transactionService.ProcessPayment(req)
	if err != nil {
		return
	}
	return
}