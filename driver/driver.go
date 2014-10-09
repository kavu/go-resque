package driver

type Enqueuer interface {
	SetClient(interface{})
	ListPush(queue string, jobJSON string) (listLength int64, err error)
}
