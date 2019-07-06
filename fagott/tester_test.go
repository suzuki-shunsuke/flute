package fagott

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_makeMsg(t *testing.T) {
	data := []struct {
		title   string
		msg     string
		srv     string
		reqName string
		exp     string
	}{
		{
			title:   "normal",
			msg:     "message",
			srv:     "service name",
			reqName: "create a user",
			exp: `message
service: service name
request name: create a user`,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			require.Equal(t, d.exp, makeMsg(d.msg, d.srv, d.reqName))
		})
	}
}

func Test_testBodyString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`"foo"`)),
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					BodyString: `"foo"`,
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testBodyString(t, d.req, d.service, d.route)
		})
	}
}

func Test_testPath(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				URL: &url.URL{
					Path: "/foo",
				},
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					Path: "/foo",
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testPath(t, d.req, d.service, d.route)
		})
	}
}

func Test_testMethod(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Method: "PUT",
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					Method: "put",
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testMethod(t, d.req, d.service, d.route)
		})
	}
}

func Test_testBodyJSON(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`[{"foo":"bar"}]`)),
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					BodyJSON: []map[string]string{
						{
							"foo": "bar",
						},
					},
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testBodyJSON(t, d.req, d.service, d.route)
		})
	}
}

func Test_testBodyJSONString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`[{"foo":"bar"}]`)),
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					BodyJSONString: `[
					{"foo":"bar"}
					]`,
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testBodyJSONString(t, d.req, d.service, d.route)
		})
	}
}

func Test_testHeader(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service *Service
		route   *Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
					"BAR": []string{"bar"},
				},
			},
			service: &Service{},
			route: &Route{
				Tester: &Tester{
					Header: http.Header{
						"FOO": []string{"foo"},
						"BAR": nil,
					},
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			testHeader(t, d.req, d.service, d.route)
		})
	}
}
