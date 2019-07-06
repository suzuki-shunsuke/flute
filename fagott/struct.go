package fagott

import (
	"net/http"
	"testing"
)

type (
	Transport struct {
		server    *Server
		Transport http.RoundTripper
	}

	Server struct {
		Services []Service
		T        *testing.T
	}

	Service struct {
		Endpoint string
		Routes   []Route
	}

	Route struct {
		Name     string
		Matcher  *Matcher
		Tester   *Tester
		Response *Response
	}

	Matcher struct {
		Match      func(req *http.Request) (bool, error)
		Method     string
		Path       string
		BodyString string
		BodyJSON   interface{}
		Header     http.Header
	}

	Tester struct {
		Test       func(*testing.T, *http.Request, *Service, *Route) (bool, error)
		Path       string
		Method     string
		BodyString string
		BodyJSON   interface{}
		Header     http.Header
	}

	Response struct {
		Base       http.Response
		Response   func(req *http.Request) (*http.Response, error)
		StatusCode int
		BodyJSON   interface{}
		BodyString string
		Header     http.Header
	}
)
