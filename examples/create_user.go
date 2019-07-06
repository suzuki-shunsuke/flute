package examples

import (
	"encoding/json"
	"net/http"
	"strings"
)

type (
	// Client is a simple API client
	Client struct {
		ClientFn func() (*http.Client, error)
		Token    string
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
	c := &http.Client{}
	if client.ClientFn != nil {
		var err error
		c, err = client.ClientFn()
		if err != nil {
			return nil, nil, err
		}
	}
	b, err := json.Marshal(user)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", "http://example.com/users", strings.NewReader(string(b)))
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
