package storage

// Storager interface for redis to make life easier
type Storager interface {
	NewOrder(string) (string, error)
	CloseOrder(string) error
	LookUp(string) bool
}
