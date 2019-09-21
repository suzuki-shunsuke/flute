package flute

import (
	"net/http"
	"net/url"
	"testing"
)

type (
	// Transport implements http.RoundTripper.
	Transport struct {
		// Each service's endpoint should be unique.
		Services []Service
		// If *testing.T is nil, the transport is a just mock and doesn't run the test.
		T *testing.T
		// Transport is used when the request doesn't match with any services.
		Transport http.RoundTripper
	}

	// Service is a service.
	Service struct {
		// The format of Endpoint should be "scheme://host", and other parameters
		// such as path and queries shouldn't be set.
		// These parameters should be set at the matcher or tester.
		Endpoint string
		// If the request matches with a route, other routes are ignored.
		Routes []Route
	}

	// Route is the pair of the macher, tester, and response.
	Route struct {
		// Name is embedded the assertion and useful to specify where the test fails.
		Name string
		// SPEC if the matcher is nil, the route matches the request.
		Matcher *Matcher
		// SPEC if the tester is nil, no test is run
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
		// PartOfQuery is the request query parameters.
		PartOfQuery url.Values
		// Query is the request query parameters.
		Query url.Values
		// BodyString is the request body.
		BodyString string
		// BodyJSON is marshaled to JSON and compared to the request body as JSON.
		BodyJSON interface{}
		// BodyJSONString is a JSON string and compared to the request body as JSON.
		BodyJSONString string
		// PartOfHeader is the request header's conditions.
		// If the header value is nil, RoundTrip checks whether the key is included in the request header.
		// Otherwise, RoundTrip also checks whether the value is equal.
		PartOfHeader http.Header
		// Header is the request header's conditions.
		Header http.Header
	}

	// Tester has the request's tests.
	Tester struct {
		Test func(*testing.T, *http.Request, *Service, *Route)
		// Path is the request path.
		Path string
		// Path is the request method such as "GET".
		Method string
		// BodyString is the request body.
		BodyString string
		// BodyJSON is marshaled to JSON and compared to the request body as JSON.
		BodyJSON interface{}
		// BodyJSONString is a JSON string and compared to the request body as JSON.
		BodyJSONString string
		// PartOfHeader is the request header's conditions.
		// If the header value is nil, RoundTrip checks whether the key is included in the request header.
		// Otherwise, RoundTrip also checks whether the value is equal.
		PartOfHeader http.Header
		// Header is the request header's conditions.
		Header http.Header
		// PartOfQuery is the request query parameters.
		// If the query value is nil, RoundTrip checks whether the key is included in the request query.
		// Otherwise, RoundTrip also checks whether the value is equal.
		PartOfQuery url.Values
		// Query is the request query parameters.
		Query url.Values
	}

	// Response has the response parameters.
	Response struct {
		// Base is the base response.
		Base http.Response
		// If Response isn't nil, Response is called to return the response and other parameters are ignored.
		Response func(req *http.Request) (*http.Response, error)
		// BodyJSON is marshaled to JSON and used as the response body.
		// BodyJSON and BodyString should only be set to one or the other.
		BodyJSON interface{}
		// BodyString is the response body.
		// BodyJSON and BodyString should only be set to one or the other.
		BodyString string
	}
)
