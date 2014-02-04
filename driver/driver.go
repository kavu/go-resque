package driver

type Enqueuer interface {
	SetClient(interface{})
	ListPush(queue string, jobJson string) (listLength int64, err error)
}
