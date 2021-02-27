package debrid

// Client represents a debrid client.
// It has the generic interface of Adapter and is backed by a service-specific client.
type Client struct {
	Adapter
}

// NewClient returns a new debrid client, backed by a service-specific client.
func NewClient(adapter Adapter) *Client {
	return &Client{adapter}
}
