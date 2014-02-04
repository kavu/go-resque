package resque

import (
	"encoding/json"
	// "github.com/davecgh/go-spew/spew"
	"github.com/kavu/go-resque/driver"
)

var drivers = make(map[string]driver.Enqueuer)

type jobArg interface{}

type job struct {
	Class string   `json:"class"`
	Args  []jobArg `json:"args"`
}

type redisEnqueuer struct {
	drv driver.Enqueuer
}

func Register(name string, driver driver.Enqueuer) {
	if _, d := drivers[name]; d {
		panic("Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func NewRedisEnqueuer(drvName string, client interface{}) *redisEnqueuer {
	drv, ok := drivers[drvName]
	if !ok {
		panic("No driver " + drvName)
	}

	drv.SetClient(client)
	return &redisEnqueuer{drv: drv}
}

func (enqueuer *redisEnqueuer) Enqueue(queue, jobClass string, args ...jobArg) (int64, error) {
	// NOTE: Dirty hack to make a [{}] JSON struct
	if len(args) == 0 {
		args = append(make([]jobArg, 0), make(map[string]jobArg, 0))
	}

	jobJson, err := json.Marshal(&job{jobClass, args})
	if err != nil {
		return -1, err
	}

	return enqueuer.drv.ListPush(queue, string(jobJson))
}
