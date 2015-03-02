package resque

import (
	"github.com/kavu/go-resque"
	"github.com/kavu/go-resque/driver"
	"github.com/vmihailenco/redis/v2"
)

func init() {
	resque.Register("redisv2", &drv{})
}

type drv struct {
	client *redis.Client
	driver.Enqueuer
}

func (d *drv) SetClient(client interface{}) {
	d.client = client.(*redis.Client)
}

func (d *drv) ListPush(queue string, jobJSON string) (int64, error) {
	return d.client.RPush(queue, jobJSON).Result()
}
