package publisher

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Publisher struct {
	Namespace string
	Pool      *redis.Pool
	enqueuer  *work.Enqueuer
}

func NewPublisher(namespace string, pool *redis.Pool) *Publisher {
	publisher := new(Publisher)
	publisher.Namespace = namespace
	publisher.Pool = pool

	return publisher
}

func (c *Publisher) Initialize() {
	c.enqueuer = work.NewEnqueuer(c.Namespace, c.Pool)
}

func (c *Publisher) Enqueue(jobName string, args map[string]interface{}) (err error) {
	_, err = c.enqueuer.Enqueue(jobName, args)
	return
}

func (c *Publisher) ScheduledEnqueue(jobName string, secondsInFuture int64, args map[string]interface{}) (err error) {
	_, err = c.enqueuer.EnqueueIn(jobName, secondsInFuture, args)
	return
}
