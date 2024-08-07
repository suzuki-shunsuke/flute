package examples //nolint:testpackage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/flute/v2/flute"
)

func TestClient_CreateUser(t *testing.T) {
	token := "XXXXX"
	client := &Client{
		Token: token,
		HTTPClient: &http.Client{
			Transport: flute.Transport{
				T: t,
				Services: []flute.Service{
					{
						Endpoint: "http://example.com",
						Routes: []flute.Route{
							{
								Name: "create a user",
								Matcher: flute.Matcher{
									Method: http.MethodPost,
									Path:   "/users",
								},
								Tester: flute.Tester{
									BodyJSONString: `{
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
									Header: http.Header{
										"Authorization": []string{"token " + token},
									},
								},
								Response: flute.Response{
									Base: http.Response{
										StatusCode: http.StatusCreated,
									},
									BodyString: `{
										  "id": 10,
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
								},
							},
						},
					},
				},
			},
		},
	}
	user, _, err := client.CreateUser(&User{ //nolint:bodyclose
		Name:  "foo",
		Email: "foo@example.com",
	})
	require.NoError(t, err)
	require.Equal(t, &User{
		ID:    10,
		Name:  "foo",
		Email: "foo@example.com",
	}, user)
}
