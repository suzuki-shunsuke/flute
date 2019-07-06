package fagott

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func createHTTPResponse(req *http.Request, resp *Response) (*http.Response, error) {
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
				StatusCode: 500,
			}, err
		}
		body = ioutil.NopCloser(strings.NewReader(string(b)))
	}
	if resp.BodyString != "" {
		body = ioutil.NopCloser(strings.NewReader(resp.BodyString))
	}
	if len(resp.Header) != 0 {
		r.Header = resp.Header
	}
	r.Body = body
	return &r, nil
}
