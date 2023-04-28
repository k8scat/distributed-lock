package redis

type Lock interface {
	Lock() (bool, error)
	Unlock() error
}
