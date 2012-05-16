package resque

import (
  "encoding/json"
  redis "github.com/simonz05/godis"
)

type jobArg interface{}

type job struct {
  Class string   `json:"class"`
  Args  []jobArg `json:"args"`
}

func Enqueue(client *redis.Client, queue, job_class string, args ...jobArg) (int64, error) {
  var j = &job{job_class, makeJobArgs(args)}

  job_json, _ := json.Marshal(j)

  return client.Lpush("resque:queue:"+queue, job_json)
}

func makeJobArgs(args []jobArg) []jobArg {
  if len(args) == 0 {
    // NOTE: Dirty hack to make a [{}] JSON struct
    return append(make([]jobArg, 0), make(map[string]jobArg, 0))
  }

  return args
}
