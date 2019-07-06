package example

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/suzuki-shunsuke/fagott/fagott"
)

func TestFoo(t *testing.T) {
	srv := &fagott.Server{
		T: t,
		Services: []fagott.Service{
			{
				Endpoint: "http://example.com",
				Routes: []fagott.Route{
					{
						Name: "get foo",
						Matcher: &fagott.Matcher{
							Method: "GET",
						},
						Tester: &fagott.Tester{
							Path: "/foo",
						},
						Response: &fagott.Response{
							Base: http.Response{
								StatusCode: 200,
							},
							BodyString: `{"message": "hello"}`,
						},
					},
				},
			},
		},
	}
	defer func(transport http.RoundTripper) {
		http.DefaultClient.Transport = transport
	}(http.DefaultClient.Transport)
	http.DefaultClient.Transport = fagott.NewTransport(srv)
	r, err := http.Get("http://example.com/foo")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestBar(t *testing.T) {
	srv := &fagott.Server{
		T: t,
		Services: []fagott.Service{
			{
				Endpoint: "http://hello.example.com",
				Routes: []fagott.Route{
					{
						Name: "get foo",
						Matcher: &fagott.Matcher{
							Method: "GET",
						},
						Tester: &fagott.Tester{
							Path: "/foo",
						},
						Response: &fagott.Response{
							Base: http.Response{
								StatusCode: 200,
							},
							BodyString: `{"message": "hello"}`,
						},
					},
				},
			},
		},
	}

	client := &http.Client{
		Transport: fagott.NewTransport(srv),
	}
	req, err := http.NewRequest("GET", "http://hello.example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}
	r, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
