package flute_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/suzuki-shunsuke/flute/v2/flute"
)

func Example_simpleMock() {
	http.DefaultClient = &http.Client{
		Transport: flute.Transport{
			// if *testing.T isn't given, the transport is a just mock and doesn't run the test.
			// T: t,
			Services: []flute.Service{
				{
					Endpoint: "http://example.com",
					Routes: []flute.Route{
						{
							Name: "get a user",
							Matcher: flute.Matcher{
								Method: "GET",
								Path:   "/users",
								Query: url.Values{
									"id": []string{"10"},
								},
							},
							Response: flute.Response{
								Base: http.Response{
									StatusCode: 201,
								},
								BodyString: `{"id": 10, "name": "foo", "email": "foo@example.com"}`,
							},
						},
					},
				},
			},
		},
	}
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://example.com/users?id=10", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	// Output:
	// {"id": 10, "name": "foo", "email": "foo@example.com"}
}
