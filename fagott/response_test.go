package fagott

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

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
			},
			exp:  &http.Response{},
			body: `{"foo":"bar"}`,
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
