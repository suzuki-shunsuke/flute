package flute

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	invalidMarshaler struct{}
)

func (*invalidMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed to marshal JSON")
}

func Test_createHTTPResponse(t *testing.T) {
	data := []struct {
		title string
		req   *http.Request
		resp  *Response
		isErr bool
		exp   *http.Response
		body  string
	}{
		{
			title: "body json isn't nil",
			req:   &http.Request{},
			resp: &Response{
				Base: http.Response{
					Header: http.Header{
						"FOO": []string{"foo"},
					},
				},
				BodyJSON: map[string]interface{}{
					"foo": "bar",
				},
			},
			exp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
			},
			body: `{"foo":"bar"}`,
		},
		{
			title: "failed to marshal json",
			req:   &http.Request{},
			resp: &Response{
				BodyJSON: &invalidMarshaler{},
			},
			isErr: true,
		},
		{
			title: "body string isn't nil",
			req:   &http.Request{},
			resp: &Response{
				BodyString: `{"foo":"bar"}`,
			},
			exp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
			},
			body: `{"foo":"bar"}`,
		},
		{
			title: "nil request body",
			req:   &http.Request{},
			resp:  &Response{},
			exp: &http.Response{
				Body: ioutil.NopCloser(strings.NewReader("")),
			},
		},
		{
			title: "resp.Response",
			req:   &http.Request{},
			resp: &Response{
				Response: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       ioutil.NopCloser(strings.NewReader("foo")),
						StatusCode: 403,
					}, nil
				},
			},
			exp: &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader("foo")),
				StatusCode: 403,
			},
			body: "foo",
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			resp, err := createHTTPResponse(d.req, d.resp)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, resp)

			// https://golang.org/pkg/net/http/#Response
			// The http Client and Transport guarantee that Body is always
			// non-nil, even on responses without a body or responses with
			// a zero-length body.
			require.NotNil(t, resp.Body)

			require.Equal(t, d.exp.StatusCode, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			require.Equal(t, d.body, string(b))
		})
	}
}
