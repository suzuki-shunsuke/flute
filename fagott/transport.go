package fagott

import (
	"net/http"
)

func (transport *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, service := range transport.Services {
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
			if transport.T != nil {
				testRequest(transport.T, req, &service, &route)
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
