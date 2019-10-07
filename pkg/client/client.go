package client

import (
	"net/url"

	"github.com/craftyphotons/ten34/pkg/client/backend"
)

const (
	// BackendSchemeRoute53 is the scheme to use for the Route53 backend
	BackendSchemeRoute53 = "route53"
)

// Client provides a database client session
type Client struct {
	// Backend is the database backend
	Backend backend.Backend

	// URI is the uniform resource identifier of the database
	URI url.URL
}

// New creates a client for a new database session
func New(uri url.URL) (*Client, error) {
	scheme := uri.Scheme
	client := &Client{URI: uri}

	switch scheme {
	case BackendSchemeRoute53:
		client.Backend = backend.NewRoute53(uri)
	}

	err := client.Backend.Setup(uri)

	return client, err
}

// CreateDB creates a database
func (c *Client) CreateDB() error {
	return c.Backend.CreateDB(c.URI)
}

// DropDB deletes a database
func (c *Client) DropDB() error {
	return c.Backend.DropDB(c.URI)
}

// Delete deletes a key from the database
func (c *Client) Delete(key string) error {
	return c.Backend.Delete(c.URI, key)
}

// Get retrieves keys from the database
func (c *Client) Get(key string) (string, error) {
	return c.Backend.Get(c.URI, key)
}

// Put writes a key-value pair to the database
func (c *Client) Put(key, val string) error {
	return c.Backend.Put(c.URI, key, val)
}
