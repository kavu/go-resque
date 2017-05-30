package driver

import "time"

type Enqueuer interface {
	SetClient(name string, client interface{})
	ListPush(queue string, jobJSON string) (listLength int64, err error)
	ListPushDelay(t time.Time, queue string, jobJSON string) (bool, error)
	Poll()
}
