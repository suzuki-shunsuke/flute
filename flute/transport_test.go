package flute_test

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/flute/flute"
	"github.com/suzuki-shunsuke/gomic/gomic"
)

func TestTransport_RoundTrip(t *testing.T) { //nolint:funlen
	token := "XXXXX"
	data := []struct {
		title     string
		req       *http.Request
		transport flute.Transport
		isErr     bool
		exp       *http.Response
	}{
		{
			title: "normal",
			req: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/users",
				},
				Method: "POST",
				Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
				Header: http.Header{
					"Authorization": []string{"token " + token},
				},
			},
			transport: flute.Transport{
				T: t,
				Services: []flute.Service{
					{
						Endpoint: "http://example.org",
					},
					{
						Endpoint: "http://example.com",
						Routes: []flute.Route{
							{
								Matcher: flute.Matcher{
									Method: "GET",
								},
							},
							{
								Name: "create a user",
								Matcher: flute.Matcher{
									Method: "POST",
									Path:   "/users",
								},
								Tester: &flute.Tester{
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
			exp: &http.Response{
				StatusCode: 201,
			},
		},
		{
			title: "failed to match",
			req: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/users",
				},
				Method: "POST",
				Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
				Header: http.Header{
					"Authorization": []string{"token " + token},
				},
			},
			transport: flute.Transport{
				Services: []flute.Service{
					{
						Endpoint: "http://example.com",
						Routes: []flute.Route{
							{
								Matcher: flute.Matcher{
									Match: func(req *http.Request) (bool, error) {
										return false, errors.New("failed to match")
									},
								},
							},
						},
					},
				},
				Transport: flute.NewMockRoundTripper(t, gomic.DoNothing).
					SetReturnRoundTrip(&http.Response{
						StatusCode: 401,
					}, nil),
			},
			exp: &http.Response{
				StatusCode: 401,
			},
		},
		{
			title: "failed to match and transport.T isn't nil",
			req: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/users",
				},
				Method: "POST",
				Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
				Header: http.Header{
					"Authorization": []string{"token " + token},
				},
			},
			transport: flute.Transport{
				T: t,
				Services: []flute.Service{
					{
						Endpoint: "http://example.com",
						Routes: []flute.Route{
							{
								Matcher: flute.Matcher{
									Match: func(req *http.Request) (bool, error) {
										return false, errors.New("failed to match")
									},
								},
							},
						},
					},
				},
				Transport: flute.NewMockRoundTripper(t, gomic.DoNothing).
					SetReturnRoundTrip(&http.Response{
						StatusCode: 401,
					}, nil),
			},
			exp: &http.Response{
				StatusCode: 401,
			},
		},
		{
			title: "noMatchedRouteRoundTrip is called",
			req: &http.Request{
				URL: &url.URL{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/users",
				},
				Method: "POST",
				Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
				Header: http.Header{
					"Authorization": []string{"token " + token},
				},
			},
			transport: flute.Transport{
				Services: []flute.Service{
					{
						Endpoint: "http://example.com",
						Routes: []flute.Route{
							{
								Matcher: flute.Matcher{
									Match: func(req *http.Request) (bool, error) {
										return false, errors.New("failed to match")
									},
								},
							},
						},
					},
				},
			},
			exp: &http.Response{
				StatusCode: 404,
			},
		},
		{
			title: "transport.Transport is called",
			req:   &http.Request{},
			transport: flute.Transport{
				Transport: flute.NewMockRoundTripper(t, gomic.DoNothing).
					SetReturnRoundTrip(&http.Response{
						StatusCode: 401,
					}, nil),
			},
			exp: &http.Response{
				StatusCode: 401,
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			resp, err := d.transport.RoundTrip(d.req)
			if resp != nil && resp.Body != nil {
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp.StatusCode, resp.StatusCode)
		})
	}
}

func BenchmarkTransport_RoundTrip(b *testing.B) { //nolint:funlen
	token := "XXXXX"
	transport := flute.Transport{
		Services: []flute.Service{
			{
				Endpoint: "http://example.org",
			},
			{
				Endpoint: "http://example.com",
				Routes: []flute.Route{
					{
						Matcher: flute.Matcher{
							Method: "GET",
						},
					},
					{
						Name: "create a user",
						Matcher: flute.Matcher{
							Method: "POST",
							Path:   "/users",
						},
						Tester: &flute.Tester{
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
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, _ := transport.RoundTrip(&http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/users",
			},
			Method: "POST",
			Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
			Header: http.Header{
				"Authorization": []string{"token " + token},
			},
		})
		_, _ = io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}
