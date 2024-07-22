package flute

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_testRequest(t *testing.T) { //nolint:funlen
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "if the tester is nil, do nothing",
			req: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path:     "/users",
					RawQuery: "name=foo&age=10",
				},
				Body: io.NopCloser(strings.NewReader(`{
				  "name": "foo",
				  "email": "foo@example.com"
				}`)),
			},
			route: Route{},
		},
		{
			title: "body json string",
			req: &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path:     "/users",
					RawQuery: "name=foo&age=10",
				},
				Body: io.NopCloser(strings.NewReader(`{
				  "name": "foo",
				  "email": "foo@example.com"
				}`)),
				Header: http.Header{
					"Authorization": []string{"token XXXXX"},
					"Content-Type":  []string{"application/json"},
				},
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					Method: http.MethodPost,
					Path:   "/users",
					BodyJSONString: `{
					  "name": "foo",
					  "email": "foo@example.com"
					}`,
					PartOfHeader: http.Header{
						"Authorization": []string{"token XXXXX"},
					},
					Header: http.Header{
						"Authorization": []string{"token XXXXX"},
						"Content-Type":  []string{"application/json"},
					},
					PartOfQuery: url.Values{
						"name": []string{"foo"},
					},
					Query: url.Values{
						"name": []string{"foo"},
						"age":  []string{"10"},
					},
					Test: func(t *testing.T, req *http.Request, service Service, route Route) {},
				},
			},
		},
		{
			title: "body string",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`foo`)),
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyString: `foo`,
				},
			},
		},
		{
			title: "body json",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`[{"name": "foo"}]`)),
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyJSON: []map[string]interface{}{
						{
							"name": "foo",
						},
					},
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testRequest(t, d.req, d.service, d.route)
		})
	}
}

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
		d := d
		t.Run(d.title, func(t *testing.T) {
			require.Equal(t, d.exp, makeMsg(d.msg, d.srv, d.reqName))
		})
	}
}

func Test_testBodyString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`"foo"`)),
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyString: `"foo"`,
				},
			},
		},
		{
			title:   "request body is nil",
			req:     &http.Request{},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyString: "",
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testBodyString(t, d.req, d.service, d.route)
		})
	}
}

func Test_testPath(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				URL: &url.URL{
					Path: "/foo",
				},
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					Path: "/foo",
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testPath(t, d.req, d.service, d.route)
		})
	}
}

func Test_testMethod(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Method: http.MethodPut,
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					Method: "put",
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testMethod(t, d.req, d.service, d.route)
		})
	}
}

func Test_testBodyJSON(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`[{"foo":"bar"}]`)),
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyJSON: []map[string]string{
						{
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			title:   "the request body is nil",
			req:     &http.Request{},
			service: Service{},
			route: Route{
				Tester: Tester{},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testBodyJSON(t, d.req, d.service, d.route)
		})
	}
}

func Test_testBodyJSONString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`[{"foo":"bar"}]`)),
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					BodyJSONString: `[
					{"foo":"bar"}
					]`,
				},
			},
		},
		{
			title:   "the request body is nil",
			req:     &http.Request{},
			service: Service{},
			route: Route{
				Tester: Tester{},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testBodyJSONString(t, d.req, d.service, d.route)
		})
	}
}

func Test_testPartOfHeader(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
					"BAR": []string{"bar"},
				},
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					Header: http.Header{
						"FOO": []string{"foo"},
						"BAR": nil,
					},
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testPartOfHeader(t, d.req, d.service, d.route)
		})
	}
}

func Test_testPartOfQuery(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		service Service
		route   Route
	}{
		{
			title: "normal",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			service: Service{},
			route: Route{
				Tester: Tester{
					Query: url.Values{
						"name": []string{"foo"},
					},
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			testPartOfQuery(t, d.req, d.service, d.route)
		})
	}
}
