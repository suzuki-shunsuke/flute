package fagott

import (
	"net/http"
)

func NewTransport(server *Server) *Transport {
	return &Transport{
		server: server,
	}
}

func (transport *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, service := range transport.server.Services {
		if !IsMatchService(req, &service) {
			continue
		}
		for _, route := range service.Routes {
			b, err := IsMatch(req, route.Matcher)
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
			transport.Test(req, &service, &route)
			// return response
			return CreateHTTPResponse(req, route.Response)
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
