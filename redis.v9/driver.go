package resque

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/skaurus/go-resque"
	"github.com/skaurus/go-resque/driver"
)

func init() {
	resque.Register("redis.v9", &drv{})
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

func (d *drv) ListPush(ctx context.Context, queue string, jobJSON string) (int64, error) {
	return d.client.RPush(ctx, d.nameSpace+"queue:"+queue, []byte(jobJSON)).Result()
}

func (d *drv) ListPushDelay(ctx context.Context, t time.Time, queue string, jobJSON string) (bool, error) {
	_, err := d.client.ZAdd(ctx, d.nameSpace+"queue:"+queue, redis.Z{
		Score:  timeToSecondsWithNanoPrecision(t),
		Member: []byte(jobJSON),
	}).Result()
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

func (d *drv) Poll(ctx context.Context) {
	go func(d *drv) {
		for {
			now := timeToSecondsWithNanoPrecision(time.Now())
			for key := range d.schedule {
				jobs, _ := d.client.ZRangeArgs(ctx, redis.ZRangeArgs{
					Key:     key,
					ByScore: true,
					Start:   "-inf",
					Stop:    now,
					Count:   1,
				}).Result()
				if len(jobs) == 0 {
					continue
				}
				if removed, _ := d.client.ZRem(ctx, key, []byte(jobs[0])).Result(); removed > 0 {
					d.client.LPush(ctx, key, []byte(jobs[0]))
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(d)
}
