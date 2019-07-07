package fagott

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzuki-shunsuke/gomic/gomic"
)

func TestTransport_RoundTrip(t *testing.T) {
	token := "XXXXX"
	data := []struct {
		title     string
		req       *http.Request
		transport *Transport
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
			transport: &Transport{
				T: t,
				Services: []Service{
					{
						Endpoint: "http://example.org",
					},
					{
						Endpoint: "http://example.com",
						Routes: []Route{
							{
								Matcher: &Matcher{
									Method: "GET",
								},
							},
							{
								Name: "create a user",
								Matcher: &Matcher{
									Method: "POST",
									Path:   "/users",
								},
								Tester: &Tester{
									BodyJSONString: `{
										  "name": "foo",
										  "email": "foo@example.com"
										}`,
									Header: http.Header{
										"Authorization": []string{"token " + token},
									},
								},
								Response: &Response{
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
			transport: &Transport{
				Services: []Service{
					{
						Endpoint: "http://example.com",
						Routes: []Route{
							{
								Matcher: &Matcher{
									Match: func(req *http.Request) (bool, error) {
										return false, errors.New("failed to match")
									},
								},
							},
						},
					},
				},
				Transport: NewMockRoundTripper(t, gomic.DoNothing).
					SetReturnRoundTrip(&http.Response{
						StatusCode: 401,
					}, nil),
			},
			exp: &http.Response{
				StatusCode: 401,
			},
		},
		{
			title: "transport.Transport is called",
			req:   &http.Request{},
			transport: &Transport{
				Transport: NewMockRoundTripper(t, gomic.DoNothing).
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
		t.Run(d.title, func(t *testing.T) {
			resp, err := d.transport.RoundTrip(d.req)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, d.exp.StatusCode, resp.StatusCode)
		})
	}
}

func Test_makeNoMatchedRouteMsg(t *testing.T) {
	data := []struct {
		title string
		req   *http.Request
		exp   string
	}{
		{
			title: "normal",
			req: &http.Request{
				URL: &url.URL{
					Scheme:   "http",
					Host:     "example.com",
					Path:     "/users",
					RawQuery: "print=true",
				},
				Method: "POST",
				Body:   ioutil.NopCloser(strings.NewReader(`{"name": "foo", "email": "foo@example.com"}`)),
				Header: http.Header{
					"Authorization": []string{"token XXXXX"},
				},
			},
			exp: `no route matches the request.
url: http://example.com/users?print=true
method: POST
query:
  print: true
header:
  Authorization: token XXXXX
body:
{"name": "foo", "email": "foo@example.com"}`,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			require.Equal(t, d.exp, makeNoMatchedRouteMsg(t, d.req))
		})
	}
}
