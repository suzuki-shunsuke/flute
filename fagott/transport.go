package fagott

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	noMatchedRouteMsgTpl = `no route matches the request.
url: %s
method: %s
query:
%s
header:
%s
body:
%s`
)

// RoundTrip implements http.RoundTripper.
// RoundTrip traverses the matched route and run the test and returns response.
func (transport *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, service := range transport.Services {
		if !isMatchService(req, &service) {
			continue
		}
		for _, route := range service.Routes {
			b, err := isMatch(req, route.Matcher)
			if err != nil {
				if transport.T != nil {
					transport.T.Logf("failed to check whether the route matches the request: %v", err)
				} else {
					fmt.Fprintf(os.Stderr, "failed to check whether the route matches the request: %v\n", err)
				}
			}
			if !b {
				continue
			}
			// test
			if transport.T != nil {
				testRequest(transport.T, req, &service, &route)
			}
			// return response
			return createHTTPResponse(req, route.Response)
		}
	}
	// no route matches the request
	if transport.Transport != nil {
		return transport.Transport.RoundTrip(req)
	}
	return noMatchedRouteRoundTrip(transport.T, req)
}

func makeNoMatchedRouteMsg(t *testing.T, req *http.Request) string {
	query := req.URL.Query()
	qArr := make([]string, len(query))
	i := 0
	for k, v := range query {
		qArr[i] = "  " + k + ": " + strings.Join(v, ", ")
		i++
	}

	hArr := make([]string, len(req.Header))
	j := 0
	for k, v := range req.Header {
		hArr[j] = "  " + k + ": " + strings.Join(v, ", ")
		j++
	}

	body := ""
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			assert.Nil(t, err, "failed to reqd the request body")
		} else {
			body = string(b)
		}
	}
	return fmt.Sprintf(
		noMatchedRouteMsgTpl,
		req.URL.String(),
		req.Method,
		strings.Join(qArr, "\n"),
		strings.Join(hArr, "\n"),
		body,
	)
}

func noMatchedRouteRoundTrip(t *testing.T, req *http.Request) (*http.Response, error) {
	if t == nil {
		return &http.Response{
			Request:    req,
			StatusCode: 404,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "no route matches the request"}`)),
		}, nil
	}

	require.Fail(t, makeNoMatchedRouteMsg(t, req))
	return &http.Response{
		Request:    req,
		StatusCode: 404,
		Body:       ioutil.NopCloser(strings.NewReader(`{"message": "no route matches the request"}`)),
	}, nil
}
