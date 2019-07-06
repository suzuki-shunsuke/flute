package fagott

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_isMatchService(t *testing.T) {
	data := []struct {
		title    string
		scheme   string
		host     string
		endpoint string
		exp      bool
	}{
		{
			title:    "normal",
			scheme:   "http",
			host:     "example.com",
			endpoint: "http://example.com",
			exp:      true,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b := isMatchService(&http.Request{
				URL: &url.URL{
					Scheme: d.scheme,
					Host:   d.host,
				},
			}, &Service{
				Endpoint: d.endpoint,
			})
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatch(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "path doesn't match",
			req: &http.Request{
				URL: &url.URL{
					Path: "/foo",
				},
			},
			matcher: &Matcher{
				Path: "/bar",
			},
		},
		{
			title: "method doesn't match",
			req: &http.Request{
				Method: "GET",
			},
			matcher: &Matcher{
				Method: "POST",
			},
		},
		{
			title: "body string doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader("foo")),
			},
			matcher: &Matcher{
				BodyString: "hello",
			},
		},
		{
			title: "body json doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: &Matcher{
				BodyJSON: 10,
			},
		},
		{
			title: "body json string doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: &Matcher{
				BodyJSONString: `"bar"`,
			},
		},
		{
			title: "header doesn't match",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			matcher: &Matcher{
				Header: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "match function doesn't match",
			matcher: &Matcher{
				Match: func(req *http.Request) (bool, error) {
					return false, nil
				},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatch(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatchQuery(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "header value doesn't match",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: &Matcher{
				Query: url.Values{
					"name": []string{"bar"},
				},
			},
		},
		{
			title: "query isn't found",
			req: &http.Request{
				URL: &url.URL{},
			},
			matcher: &Matcher{
				Query: url.Values{
					"name": nil,
				},
			},
		},
		{
			title: "query matches",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: &Matcher{
				Query: url.Values{
					"name": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatchQuery(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatchHeader(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "header value doesn't match",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			matcher: &Matcher{
				Header: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "header isn't found (nil)",
			req: &http.Request{
				Header: http.Header{},
			},
			matcher: &Matcher{
				Header: http.Header{
					"FOO": nil,
				},
			},
		},
		{
			title: "header matches",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			matcher: &Matcher{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatchHeader(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatchBodyString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body matches",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader("foo")),
			},
			matcher: &Matcher{
				BodyString: "foo",
			},
			exp: true,
		},
		{
			title: "request body doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader("foo")),
			},
			matcher: &Matcher{
				BodyString: "bar",
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatchBodyString(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatchBodyJSONString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: &Matcher{
				BodyJSONString: `{"name": "foo", "id": 10}`,
			},
			exp: true,
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: &Matcher{
				BodyJSONString: `{"name": "foo", "id": 9}`,
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatchBodyJSONString(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Test_isMatchBodyJSON(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher *Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`{"name": "foo"}`)),
			},
			matcher: &Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
			exp: true,
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: ioutil.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: &Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
		},
	}

	for _, d := range data {
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatchBodyJSON(d.req, d.matcher)
			if d.isErr {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}
