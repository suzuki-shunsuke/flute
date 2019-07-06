# fagott

Golang HTTP Client testing framework

## Overview

`fagott` is the Golang HTTP client testing framework.
The goal is

* Test request parameters such as the request path, headers and body
* Mock the HTTP server

We assume `fagott` is used at test functions.

`fagott.NewTransport()` returns `*fagott.Transport`, which implements [http.RoundTripper](https://golang.org/pkg/net/http/#RoundTripper).
So please set the returned value to `http.Client.Transport`.

`fagott` uses [testify](https://github.com/stretchr/testify)'s assert internally.
You can test the http request parameters with assert.

For example, The following test fails because the required request header "FOO" isn't set.

```go
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
							Header: http.Header{
								"FOO": nil,
							},
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
```

```console
$ go test -v ./... -covermode=atomic
=== RUN   TestBar
{"message": "hello"}
--- FAIL: TestBar (0.00s)
    transport.go:175:
                Error Trace:    transport.go:175
                                                        transport.go:72
                Error:          the request header should set: FOO
                                service: http://hello.example.com
                                request name: get foo
                Test:           TestBar
```

## How to use


```go
client := &http.Client{
	Transport: fagott.NewTransport(
		&fagott.Server{
			T: t,
			Services: []fagott.Service{
				{
					Endpoint: "http://hello.example.com",
					Routes: []fagott.Route{
						{
							Name: "get foo",
							Matcher: &fagott.Matcher{
								Method: "GET",
								Path: "/foo",
							},
							Tester: &fagott.Tester{
								Header: http.Header{
									"FOO": nil,
								},
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
		}),
}
```

```go
defer func(transport http.RoundTripper) {
	http.DefaultClient.Transport = transport
}(http.DefaultClient.Transport)
http.DefaultClient.Transport = fagott.NewTransport(
	&fagott.Server{
		T: t,
		Services: []fagott.Service{
			{
				Endpoint: "http://hello.example.com",
				Routes: []fagott.Route{
					{
						Name: "get foo",
						Matcher: &fagott.Matcher{
							Method: "GET",
							Path: "/foo",
						},
						Tester: &fagott.Tester{
							Header: http.Header{
								"FOO": nil,
							},
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
	})
```

## License

[MIT](LICENSE)
