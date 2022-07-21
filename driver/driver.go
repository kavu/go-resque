package driver

import (
	"context"
	"time"
)

type Enqueuer interface {
	SetClient(name string, client interface{})
	ListPush(ctx context.Context, queue string, jobJSON string) (listLength int64, err error)
	ListPushDelay(ctx context.Context, t time.Time, queue string, jobJSON string) (bool, error)
	Poll(ctx context.Context)
}
