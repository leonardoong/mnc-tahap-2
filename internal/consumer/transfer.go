package consumer

import (
	"github.com/gocraft/work"
	"github.com/leonardoong/e-wallet/internal/domain/entity"
	"github.com/leonardoong/e-wallet/internal/service"
)

type transferWorker struct {
	transactionService service.ITransactionService
	workerPool         *work.WorkerPool
	jobName            string
}

func newTransferWorker(srv service.ITransactionService, pool *work.WorkerPool) *transferWorker {
	return &transferWorker{
		transactionService: srv,
		workerPool:         pool,
	}
}

func (c *transferWorker) runTransferConsumer(maxFails uint) {
	c.workerPool.JobWithOptions(c.jobName, work.JobOptions{MaxFails: maxFails}, c.processTransfer)
}

func (c *transferWorker) processTransfer(job *work.Job) (err error) {
	if err = job.ArgError(); err != nil {
		return err
	}

	req := entity.TransferRequest{
		TransferID: job.ArgString("transfer_id"),
		TargetTransferID: job.ArgString("target_transfer_id"),
		Amount:  job.ArgFloat64("amount"),
		UserID:  job.ArgString("user_id"),
		TargetUser:  job.ArgString("target_user"),
		Remarks:  job.ArgString("remarks"),
	}
	err = c.transactionService.ProcessTransfer(req)
	if err != nil {
		return
	}
	return
}
