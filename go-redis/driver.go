package resque

import (
	"github.com/fiorix/go-redis/redis"
	"github.com/kavu/go-resque"
	"github.com/kavu/go-resque/driver"
)

func init() {
	resque.Register("redis-go", &drv{})
}

type drv struct {
	client *redis.Client
	driver.Enqueuer
}

func (d *drv) SetClient(client interface{}) {
	d.client = client.(*redis.Client)
}

func (d *drv) ListPush(queue string, jobJson string) (int64, error) {
	ret_int, err := d.client.LPush(queue, jobJson)
	ret_int64 := int64(ret_int)
	return ret_int64, err
}
