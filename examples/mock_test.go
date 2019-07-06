package examples

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/suzuki-shunsuke/fagott/fagott"
)

func Example_simpleMock() {
	http.DefaultClient = &http.Client{
		Transport: &fagott.Transport{
			// if *testing.T isn't given, the transport is a just mock and doesn't run the test.
			// T: t,
			Services: []fagott.Service{
				{
					Endpoint: "http://example.com",
					Routes: []fagott.Route{
						{
							Name: "get a user",
							Matcher: &fagott.Matcher{
								Method: "GET",
								Path:   "/users/10",
							},
							Response: &fagott.Response{
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
	resp, err := http.Get("http://example.com/users/10")
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
