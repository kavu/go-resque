package resque

import (
	"strconv"
	"strings"
	"time"

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
	schedule  map[string]struct{}
	nameSpace string
}

func (d *drv) SetClient(name string, client interface{}) {
	d.client = client.(*redis.Conn)
	d.schedule = make(map[string]struct{})
	d.nameSpace = name
}

func (d *drv) ListPush(queue string, jobJSON string) (int64, error) {
	resp, err := (*d.client).Do("RPUSH", d.nameSpace+"queue:"+queue, jobJSON)
	if err != nil {
		return -1, err
	}

	return redis.Int64(resp, err)
}
func (d *drv) ListPushDelay(t time.Time, queue string, jobJSON string) (bool, error) {
	_, err := (*d.client).Do("ZADD", queue, timeToSecondsWithNanoPrecision(t), jobJSON)
	if err != nil {
		return false, err
	}
	if _, ok := d.schedule[queue]; !ok {
		d.schedule[queue] = struct{}{}
	}
	return true, nil
}
func timeToSecondsWithNanoPrecision(t time.Time) float64 {
	return float64(t.UnixNano()) / 1000000000.0 //nanoSecondPrecision
}

func (d *drv) Poll() {
	go func(d *drv) {
		for {
			for key := range d.schedule {
				now := timeToSecondsWithNanoPrecision(time.Now())
				jobs, _ := redis.Strings((*d.client).Do("ZRANGEBYSCORE", key, "-inf",
					strconv.FormatFloat(now, 'E', -1, 64)))
				if len(jobs) == 0 {
					continue
				}
				if _, err := (*d.client).Do("ZREM", key, jobs[0]); err != nil {
					queue := strings.TrimPrefix(key, d.nameSpace)
					(*d.client).Do("LPUSH", d.nameSpace+"queue:"+queue, jobs[0])
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}(d)
}
