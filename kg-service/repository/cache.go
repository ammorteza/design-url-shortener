package repository

type Cache interface {
	Ping() error
	Set(key, value string) error
	Get(key string) (string, error)
	Remove(key string) error
	Push(listName, val string) error
	Pop(listName string) (val string, err error)
	ListLen(listName string) (int64, error)
}