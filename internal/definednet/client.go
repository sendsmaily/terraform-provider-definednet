package definednet

// NewClient creates a Defined.net HTTP API client.
func NewClient(endpoint, token string) Client {
	return &client{
		endpoint: endpoint,
		token:    token,
	}
}

// Client is a Defined.net HTTP API client.
type Client interface{}

type client struct {
	endpoint string
	token    string
}

var _ Client = (*client)(nil)
