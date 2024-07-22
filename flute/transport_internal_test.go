package flute

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
				Method: http.MethodPost,
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
		d := d
		t.Run(d.title, func(t *testing.T) {
			require.Equal(t, d.exp, makeNoMatchedRouteMsg(t, d.req))
		})
	}
}

func Test_noMatchedRouteRoundTrip(t *testing.T) {
	data := []struct {
		t          *testing.T
		title      string
		req        *http.Request
		statusCode int
		isErr      bool
	}{
		{
			t:     nil,
			title: "normal",
			req: &http.Request{
				URL: &url.URL{},
			},
			statusCode: 404,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			resp, err := noMatchedRouteRoundTrip(d.t, d.req)
			if resp != nil && resp.Body != nil {
				_, _ = io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
			if d.isErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			require.Equal(t, d.statusCode, resp.StatusCode)
		})
	}
}
