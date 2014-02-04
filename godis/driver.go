package resque

import (
	"github.com/kavu/go-resque"
	"github.com/kavu/go-resque/driver"
	"github.com/simonz05/godis/redis"
)

func init() {
	resque.Register("godis", &drv{})
}

type drv struct {
	client *redis.Client
	driver.Enqueuer
}

func (d *drv) SetClient(client interface{}) {
	d.client = client.(*redis.Client)
}

func (d *drv) ListPush(queue string, jobJson string) (int64, error) {
	return d.client.Lpush(queue, jobJson)
}
