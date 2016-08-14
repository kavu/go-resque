package resque

import (
	"fmt"
	"strings"
	"time"

	"github.com/fiorix/go-redis/redis"
	"github.com/jazibjohar/go-resque"
	"github.com/jazibjohar/go-resque/driver"
)

func init() {
	resque.Register("redis-go", &drv{})
}

type drv struct {
	client *redis.Client
	driver.Enqueuer
	schedule  map[string]struct{}
	nameSpace string
}

func (d *drv) SetClient(name string, client interface{}) {
	d.client = client.(*redis.Client)
	d.schedule = make(map[string]struct{})
	d.nameSpace = name
}

func (d *drv) ListPush(queue string, jobJSON string) (int64, error) {
	listLength, err := d.client.RPush(d.nameSpace+"queue:"+queue, jobJSON)
	if err != nil {
		return -1, err
	}

	return int64(listLength), err
}
func (d *drv) ListPushDelay(t time.Time, queue string, jobJSON string) (bool, error) {
	_, err := d.client.ZAdd(queue, t.UnixNano(), jobJSON)
	if err != nil {
		return false, err
	}
	if _, ok := d.schedule[queue]; !ok {
		d.schedule[queue] = struct{}{}
	}
	return true, nil
}

func (d *drv) Poll() {
	go func(d *drv) {
		for {
			for key := range d.schedule {
				now := time.Now()
				k := fmt.Sprintf("%s -inf %d", key, now.UnixNano())
				jobs, _ := d.client.ZRangeByScore(k, 0, 1, true, true, 0, 1)
				if len(jobs) == 0 {
					continue
				}
				removed, _ := d.client.ZRem(key, jobs[0])
				if removed == 0 {
					queue := strings.TrimPrefix(key, d.nameSpace)
					d.client.LPush(d.nameSpace+"queue:"+queue, jobs[0])
				}
			}
			time.Sleep(1 * time.Second)
		}
	}(d)
}
