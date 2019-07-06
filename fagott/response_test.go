package fagott

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
				BodyJSON: map[string]interface{}{
					"foo": "bar",
				},
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			exp:  &http.Response{},
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
			exp:  &http.Response{},
			body: `{"foo":"bar"}`,
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
			require.Equal(t, d.exp.StatusCode, resp.StatusCode)
			b, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)
			require.Equal(t, d.body, string(b))
		})
	}
}
