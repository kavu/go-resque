package resque

import (
	"encoding/json"
	"time"

	"github.com/jazibjohar/go-resque/driver"
)

var drivers = make(map[string]driver.Enqueuer)

type jobArg interface{}

type MockRedisDriver struct {
	driver.Enqueuer
}

type job struct {
	Queue string   `json:"queue,omitempty"`
	Class string   `json:"class"`
	Args  []jobArg `json:"args"`
}

func Register(name string, driver driver.Enqueuer) {
	if _, d := drivers[name]; d {
		panic("Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func NewRedisEnqueuer(drvName string, client interface{}, nameSpace string) *RedisEnqueuer {
	drv, ok := drivers[drvName]
	if !ok {
		panic("No such driver: " + drvName)
	}

	drv.SetClient(nameSpace, client)
	drv.Poll()
	return &RedisEnqueuer{drv: drv}
}

type RedisEnqueuer struct {
	drv driver.Enqueuer
}

func (enqueuer *RedisEnqueuer) Enqueue(queue, jobClass string, args ...jobArg) (int64, error) {
	// NOTE: Dirty hack to make a [{}] JSON struct
	if len(args) == 0 {
		args = append(make([]jobArg, 0), make(map[string]jobArg, 0))
	}

	jobJSON, err := json.Marshal(&job{Class: jobClass, Args: args})
	if err != nil {
		return -1, err
	}

	return enqueuer.drv.ListPush(queue, string(jobJSON))
}

// EnqueueIn enque a job at a duration
func (enqueuer *RedisEnqueuer) EnqueueIn(delay time.Duration, queue, jobClass string, args ...jobArg) (bool, error) {
	enqueueTime := time.Now().Add(delay)

	if len(args) == 0 {
		args = append(make([]jobArg, 0), make(map[string]jobArg, 0))
	}

	jobJSON, err := json.Marshal(&job{Class: jobClass, Args: args, Queue: queue})
	if err != nil {
		return false, err
	}
	return enqueuer.drv.ListPushDelay(enqueueTime, queue, string(jobJSON))
}
