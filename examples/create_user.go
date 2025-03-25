package examples

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type (
	// Client is a simple API client
	Client struct {
		HTTPClient *http.Client
		Token      string
	}

	// User is a User
	User struct {
		ID    int    `json:"id,omitempty"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

// CreateUser create a user.
func (client *Client) CreateUser(user *User) (*User, *http.Response, error) {
	c := client.HTTPClient
	if c == nil {
		c = http.DefaultClient
	}
	b, err := json.Marshal(user)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com/users", strings.NewReader(string(b)))
	if err != nil {
		return nil, nil, err
	}
	req.Header = http.Header{
		"Authorization": []string{"token " + client.Token},
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	u := &User{}
	return u, resp, json.NewDecoder(resp.Body).Decode(u)
}
