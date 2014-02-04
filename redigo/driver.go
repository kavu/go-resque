package resque

import (
	"github.com/garyburd/redigo/redis"
	"github.com/kavu/go-resque"
	"github.com/kavu/go-resque/driver"
)

func init() {
	resque.Register("redigo", &drv{})
}

type drv struct {
	client *redis.Conn
	driver.Enqueuer
}

func (d *drv) SetClient(client interface{}) {
	d.client = client.(*redis.Conn)
}

func (d *drv) ListPush(queue string, jobJson string) (int64, error) {
	resp, err := (*d.client).Do("LPUSH", queue, jobJson)
	if err != nil {
		return -1, err
	}

	return redis.Int64(resp, err)
}
