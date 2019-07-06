package fagott

import (
	"net/http"
	"testing"
)

func NewTransport(t *testing.T, server *Server) *Transport {
	return &Transport{
		server: server,
		t:      t,
	}
}

func (transport *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, service := range transport.server.Services {
		if !isMatchService(req, &service) {
			continue
		}
		for _, route := range service.Routes {
			b, err := isMatch(req, route.Matcher)
			if err != nil {
				return &http.Response{
					Request:    req,
					StatusCode: 500,
				}, err
			}
			if !b {
				continue
			}
			// test
			if transport.t != nil {
				testRequest(transport.t, req, &service, &route)
			}
			// return response
			return createHTTPResponse(req, route.Response)
		}
	}
	// there is no match response
	if transport.Transport != nil {
		return transport.Transport.RoundTrip(req)
	}
	if http.DefaultClient.Transport != transport {
		return http.DefaultClient.Transport.RoundTrip(req)
	}
	return http.DefaultTransport.RoundTrip(req)
}
