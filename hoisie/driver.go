package resque

import (
	"strings"
	"time"

	"github.com/hoisie/redis"
	"github.com/jazibjohar/go-resque"
	"github.com/jazibjohar/go-resque/driver"
)

func init() {
	resque.Register("hoisie", &drv{})
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
	err := d.client.Rpush(d.nameSpace+"queue:"+queue, []byte(jobJSON))
	if err != nil {
		return -1, err
	}

	listLength, err := d.client.Llen(queue)

	return int64(listLength), err
}

func (d *drv) ListPushDelay(t time.Time, queue string, jobJSON string) (bool, error) {
	_, err := d.client.Zadd(queue, []byte(jobJSON), timeToSecondsWithNanoPrecision(t))
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
				r, _ := d.client.Zrangebyscore(key+"-inf",
					now, 1)
				var jobs []string
				for _, job := range r {
					jobs = append(jobs, string(job))
				}
				if len(jobs) == 0 {
					continue
				}
				if removed, _ := d.client.Zrem(key, []byte(jobs[0])); removed {
					queue := strings.TrimPrefix(key, d.nameSpace)
					d.client.Lpush(d.nameSpace+"queue:"+queue, []byte(jobs[0]))
				}
			}
		}
	}(d)
}
