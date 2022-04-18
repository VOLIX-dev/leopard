package caching

// Driver is the caching driver interface.
type Driver interface {

	// Get retrieves the value for the given key.
	Get(key string, target any) (bool, error)

	// Set sets the value for the given key.
	Set(key string, value any) error

	// SetTTL sets the value for the given key with a TTL.
	SetTTL(key string, value any, ttl int) error

	// Delete deletes the value for the given key.
	Delete(key string) error

	// Close closes the driver.
	Close() error

	// Open opens the driver.
	Open() error
}
