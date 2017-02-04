package mongo_test

import (
	"os"
	"testing"
	"time"

	"github.com/blankrobot/pulpe/mock"
	"github.com/blankrobot/pulpe/mongo"
)

const (
	defaultURI = "mongodb://localhost:27017/pulpe-tests"
)

// Client is a test wrapper for mongo.Client.
type Client struct {
	*mongo.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	// Create client wrapper.
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = defaultURI
	}

	c := &Client{
		Client: mongo.NewClient(uri),
	}
	c.Now = func() time.Time { return mock.Now }

	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient(t *testing.T) *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		t.Error(err)
	}

	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	err := c.Session.DB("").DropDatabase()
	if err != nil {
		return err
	}

	return c.Client.Close()
}
