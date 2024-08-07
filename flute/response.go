package flute

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func createHTTPResponse(req *http.Request, resp Response) (*http.Response, error) {
	if resp.Response != nil {
		return resp.Response(req)
	}
	r := resp.Base
	r.Request = req
	var body io.ReadCloser
	if resp.BodyJSON != nil {
		b, err := json.Marshal(resp.BodyJSON)
		if err != nil {
			return &http.Response{
				Request:    req,
				StatusCode: http.StatusInternalServerError,
			}, err
		}
		body = io.NopCloser(strings.NewReader(string(b)))
	}
	if resp.BodyString != "" {
		body = io.NopCloser(strings.NewReader(resp.BodyString))
	}
	if body == nil {
		// https://golang.org/pkg/net/http/#Response
		// The http Client and Transport guarantee that Body is always
		// non-nil, even on responses without a body or responses with
		// a zero-length body. It is the caller's responsibility to
		// close Body.
		body = io.NopCloser(strings.NewReader(""))
	}

	r.Body = body
	return &r, nil
}
