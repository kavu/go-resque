package resque

import (
	"github.com/hoisie/redis"
	"github.com/kavu/go-resque"
	"github.com/kavu/go-resque/driver"
)

func init() {
	resque.Register("hoisie", &drv{})
}

type drv struct {
	client *redis.Client
	driver.Enqueuer
}

func (d *drv) SetClient(client interface{}) {
	d.client = client.(*redis.Client)
}

func (d *drv) ListPush(queue string, jobJson string) (int64, error) {
	err := d.client.Lpush(queue, []byte(jobJson))
	if err != nil {
		return -1, err
	}

	listLength, err := d.client.Llen(queue)

	return int64(listLength), err
}
