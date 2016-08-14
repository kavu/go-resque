package resque

import (
	"strconv"
	"strings"
	"time"

	"github.com/jazibjohar/go-resque"
	"github.com/jazibjohar/go-resque/driver"
	"github.com/simonz05/godis/redis"
)

func init() {
	resque.Register("godis", &drv{})
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
	return d.client.Rpush(queue, jobJSON)
}

func (d *drv) ListPushDelay(t time.Time, queue string, jobJSON string) (bool, error) {
	_, err := d.client.Zadd(queue, timeToSecondsWithNanoPrecision(t), jobJSON)
	if err != nil {
		return false, err
	}
	if _, ok := d.schedule[queue]; !ok {
		d.schedule[queue] = struct{}{}
	}
	return true, nil
}
func timeToSecondsWithNanoPrecision(t time.Time) float64 {
	return float64(t.UnixNano()) / 1000000000.0 // nanoSecondPrecision
}

func (d *drv) Poll() {
	go func(d *drv) {
		for {
			for key := range d.schedule {
				now := timeToSecondsWithNanoPrecision(time.Now())
				r, _ := d.client.Zrangebyscore(key, "-inf",
					strconv.FormatFloat(now, 'E', -1, 64))
				jobs := r.StringArray()
				if len(jobs) == 0 {
					continue
				}
				if removed, _ := d.client.Zrem(key, jobs[0]); removed {
					queue := strings.TrimPrefix(key, d.nameSpace)
					d.client.Lpush(d.nameSpace+"queue:"+queue, jobs[0])
				}
			}
			time.Sleep(1 * time.Second)
		}
	}(d)
}
