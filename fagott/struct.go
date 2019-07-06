package fagott

import (
	"net/http"
	"testing"
)

type (
	// Transport implements http.RoundTripper.
	Transport struct {
		server *Server
		// Transport is used when the request doesn't match with any services.
		Transport http.RoundTripper
	}

	Server struct {
		Services []Service
		T        *testing.T
	}

	Service struct {
		// The format of Endpoint should be "scheme://host", and other parameters
		// such as path and queries shouldn't be set.
		// These parameters should be set at the matcher or tester.
		Endpoint string
		Routes   []Route
	}

	// Route is the pair of the macher, tester, and response.
	Route struct {
		// Name is embedded the assertion and useful to specify where the test fails.
		Name     string
		Matcher  *Matcher
		Tester   *Tester
		Response *Response
	}

	// Matcher has conditions the request matches with the route.
	Matcher struct {
		// Match is a custom function to check the request matches with the route.
		Match func(req *http.Request) (bool, error)
		// Path is the request method such as "GET".
		Method string
		// Path is the request path.
		Path string
		// BodyString is the request body.
		BodyString string
		// BodyJSON is marshaled to JSON and compared to the request body as JSON.
		BodyJSON interface{}
		// Header is the request header's conditions.
		// If the header value is nil, RoundTrip checks whether the key is included in the request header.
		// Otherwise, RoundTrip also checks whether the value is equal.
		Header http.Header
	}

	// Tester has the request's tests.
	Tester struct {
		Test func(*testing.T, *http.Request, *Service, *Route) (bool, error)
		// Path is the request path.
		Path string
		// Path is the request method such as "GET".
		Method string
		// BodyString is the request body.
		BodyString string
		// BodyJSON is marshaled to JSON and compared to the request body as JSON.
		BodyJSON interface{}
		// Header is the request header's conditions.
		// If the header value is nil, RoundTrip checks whether the key is included in the request header.
		// Otherwise, RoundTrip also checks whether the value is equal.
		Header http.Header
	}

	// Response has the response parameters.
	Response struct {
		// Base is the base response.
		Base http.Response
		// If Response isn't nil, Response is called to return the response and other parameters are ignored.
		Response func(req *http.Request) (*http.Response, error)
		// StatusCode is the response's status code.
		StatusCode int
		// BodyJSON is marshaled to JSON and used as the response body.
		// BodyJSON and BodyString should only be set to one or the other.
		BodyJSON interface{}
		// BodyString is the response body.
		// BodyJSON and BodyString should only be set to one or the other.
		BodyString string
		// Header is the response headers.
		Header http.Header
	}
)
