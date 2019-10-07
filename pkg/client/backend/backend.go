package backend

import (
	"net/url"
)

// Backend provides a database backend implementation
type Backend interface {
	// Setup sets up a backend session
	Setup(db url.URL) error

	// CreateDB creates a database
	CreateDB(db url.URL) error

	// DropDB deletes a database
	DropDB(db url.URL) error

	// Delete deletes a key from the database
	Delete(db url.URL, key string) error

	// Get retrieves keys from the database
	Get(db url.URL, key string) (string, error)

	// Put writes a key-value pair to the database
	Put(db url.URL, key, val string) error
}
