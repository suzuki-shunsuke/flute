package examples

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/fagott/fagott"
)

func TestClient_CreateUser(t *testing.T) {
	token := "XXXXX"
	client := &Client{
		Token: token,
		HTTPClient: &http.Client{
			Transport: &fagott.Transport{
				T: t,
				Services: []fagott.Service{
					{
						Endpoint: "http://example.com",
						Routes: []fagott.Route{
							{
								Name: "create a user",
								Matcher: &fagott.Matcher{
									Method: "POST",
									Path:   "/users",
								},
								Tester: &fagott.Tester{
									BodyJSONString: `{
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
									Header: http.Header{
										"Authorization": []string{"token " + token},
									},
								},
								Response: &fagott.Response{
									Base: http.Response{
										StatusCode: 201,
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
	user, _, err := client.CreateUser(&User{
		Name:  "foo",
		Email: "foo@example.com",
	})
	require.Nil(t, err)
	require.Equal(t, &User{
		ID:    10,
		Name:  "foo",
		Email: "foo@example.com",
	}, user)
}
