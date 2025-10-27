package definednet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/samber/lo"
)

// NewClient creates a Defined.net HTTP API client.
func NewClient(endpoint, token string, version string) (Client, error) {
	if lo.IsEmpty(strings.TrimSpace(endpoint)) {
		return nil, errors.New("endpoint URL must be set")
	}

	endpointURL, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}

	if lo.IsEmpty(strings.TrimSpace(token)) {
		return nil, errors.New("authorization token must be set")
	}

	return &client{
		endpoint: endpointURL,
		token:    token,
		version:  version,
	}, nil
}

// Client is a Defined.net HTTP API client.
type Client interface {
	Do(ctx context.Context, method string, path []string, request, response any) error
}

// Response is a generic data model for Defined.net responses.
type Response[D any] struct {
	Data D `json:"data"`
}

type client struct {
	endpoint *url.URL
	token    string
	version  string
}

func (c *client) Do(ctx context.Context, method string, path []string, reqPayload, respPayload any) error {
	var buf bytes.Buffer

	if reqPayload != nil {
		if err := json.NewEncoder(&buf).Encode(reqPayload); err != nil {
			return fmt.Errorf("error encoding request payload: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.endpoint.JoinPath(lo.Map(path, func(p string, _ int) string {
			return url.PathEscape(p)
		})...).String(),
		&buf,
	)

	if err != nil {
		return fmt.Errorf("error compiling HTTP request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("User-Agent", fmt.Sprintf("Terraform-smaily-definednet/%s", c.version))
	if reqPayload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		buf.Reset() // Let's recycle it.
		if _, err := io.Copy(&buf, resp.Body); err != nil {
			return fmt.Errorf("error reading error response: %w", err)
		}

		return fmt.Errorf("code=%d reason=%s", resp.StatusCode, buf.String())
	}

	if respPayload != nil {
		if err := json.NewDecoder(resp.Body).Decode(respPayload); err != nil {
			return fmt.Errorf("error decoding response payload: %w", err)
		}
	}

	return nil
}
