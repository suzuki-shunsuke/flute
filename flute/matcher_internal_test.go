package flute

import (
	"io"
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
		d := d
		t.Run(d.title, func(t *testing.T) {
			b := isMatchService(&http.Request{
				URL: &url.URL{
					Scheme: d.scheme,
					Host:   d.host,
				},
			}, Service{
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

func Benchmark_isMatchService(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isMatchService(&http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "example.com",
			},
		}, Service{
			Endpoint: "http://example.com",
		})
	}
}

func Test_isMatch(t *testing.T) { //nolint:funlen
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
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
			matcher: Matcher{
				Path: "/bar",
			},
		},
		{
			title: "method doesn't match",
			req: &http.Request{
				Method: http.MethodGet,
			},
			matcher: Matcher{
				Method: http.MethodPost,
			},
		},
		{
			title: "body string doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "hello",
			},
		},
		{
			title: "body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: Matcher{
				BodyJSON: 10,
			},
		},
		{
			title: "body json string doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: Matcher{
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
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "header isn't equal",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
					"BAR": []string{"bar"},
				},
			},
			matcher: Matcher{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
		},
		{
			title: "query doesn't match",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"bar"},
				},
			},
		},
		{
			title: "query isn't equal",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo&age=10",
				},
			},
			matcher: Matcher{
				Query: url.Values{
					"name": []string{"foo"},
				},
			},
		},
		{
			title: "match function doesn't match",
			matcher: Matcher{
				Match: func(req *http.Request) (bool, error) {
					return false, nil
				},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := isMatch(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_isMatch(b *testing.B) { //nolint:funlen
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
	}{
		{
			title: "path doesn't match",
			req: &http.Request{
				URL: &url.URL{
					Path: "/foo",
				},
			},
			matcher: Matcher{
				Path: "/bar",
			},
		},
		{
			title: "method doesn't match",
			req: &http.Request{
				Method: http.MethodGet,
			},
			matcher: Matcher{
				Method: http.MethodPost,
			},
		},
		{
			title: "body string doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "hello",
			},
		},
		{
			title: "body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: Matcher{
				BodyJSON: 10,
			},
		},
		{
			title: "body json string doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`"foo"`)),
			},
			matcher: Matcher{
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
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "header isn't equal",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
					"BAR": []string{"bar"},
				},
			},
			matcher: Matcher{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
		},
		{
			title: "query doesn't match",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"bar"},
				},
			},
		},
		{
			title: "query isn't equal",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo&age=10",
				},
			},
			matcher: Matcher{
				Query: url.Values{
					"name": []string{"foo"},
				},
			},
		},
		{
			title: "match function doesn't match",
			matcher: Matcher{
				Match: func(req *http.Request) (bool, error) {
					return false, nil
				},
			},
		},
	}

	b.ResetTimer()
	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = isMatch(d.req, d.matcher)
			}
		})
	}
}

func Test_matchPartOfQuery(t *testing.T) { //nolint:funlen
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		exp     bool
		isErr   bool
	}{
		{
			title: "query value doesn't match",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"bar"},
				},
			},
		},
		{
			title: "query isn't found",
			req: &http.Request{
				URL: &url.URL{},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
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
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := matchPartOfQuery(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_matchPartOfQuery(b *testing.B) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		exp     bool
	}{
		{
			title: "query value doesn't match",
			req: &http.Request{
				URL: &url.URL{
					RawQuery: "name=foo",
				},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"bar"},
				},
			},
		},
		{
			title: "query isn't found",
			req: &http.Request{
				URL: &url.URL{},
			},
			matcher: Matcher{
				PartOfQuery: url.Values{
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
			matcher: Matcher{
				PartOfQuery: url.Values{
					"name": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = matchPartOfQuery(d.req, d.matcher)
			}
		})
	}
}

func Test_matchPartOfHeader(t *testing.T) { //nolint:funlen
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		exp     bool
		isErr   bool
	}{
		{
			title: "header value doesn't match",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "header isn't found",
			req: &http.Request{
				Header: http.Header{},
			},
			matcher: Matcher{
				PartOfHeader: http.Header{
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
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := matchPartOfHeader(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_matchPartOfHeader(b *testing.B) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		exp     bool
	}{
		{
			title: "header value doesn't match",
			req: &http.Request{
				Header: http.Header{
					"FOO": []string{"foo"},
				},
			},
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"bar"},
				},
			},
		},
		{
			title: "header isn't found",
			req: &http.Request{
				Header: http.Header{},
			},
			matcher: Matcher{
				PartOfHeader: http.Header{
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
			matcher: Matcher{
				PartOfHeader: http.Header{
					"FOO": []string{"foo"},
				},
			},
			exp: true,
		},
	}

	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = matchPartOfHeader(d.req, d.matcher)
			}
		})
	}
}

func Test_matchBodyString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			matcher: Matcher{
				BodyString: "foo",
			},
			req: &http.Request{},
		},
		{
			title: "request body matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "foo",
			},
			exp: true,
		},
		{
			title: "request body doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "bar",
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := matchBodyString(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_matchBodyString(b *testing.B) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "foo",
			},
		},
		{
			title: "request body doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader("foo")),
			},
			matcher: Matcher{
				BodyString: "bar",
			},
		},
	}

	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = matchBodyString(d.req, d.matcher)
			}
		})
	}
}

func Test_matchBodyJSONString(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
			matcher: Matcher{
				BodyJSONString: `{"name": "foo", "id": 10}`,
			},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSONString: `{"name": "foo", "id": 10}`,
			},
			exp: true,
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSONString: `{"name": "foo", "id": 9}`,
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := matchBodyJSONString(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_matchBodyJSONString(b *testing.B) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSONString: `{"name": "foo", "id": 10}`,
			},
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSONString: `{"name": "foo", "id": 9}`,
			},
		},
	}

	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = matchBodyJSONString(d.req, d.matcher)
			}
		})
	}
}

func Test_matchBodyJSON(t *testing.T) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
		isErr   bool
		exp     bool
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
			matcher: Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
			exp: true,
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			b, err := matchBodyJSON(d.req, d.matcher)
			if d.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if d.exp {
				require.True(t, b)
				return
			}
			require.False(t, b)
		})
	}
}

func Benchmark_matchBodyJSON(b *testing.B) {
	data := []struct {
		title   string
		req     *http.Request
		matcher Matcher
	}{
		{
			title: "request body is nil",
			req:   &http.Request{},
		},
		{
			title: "request body json matches",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
		},
		{
			title: "request body json doesn't match",
			req: &http.Request{
				Body: io.NopCloser(strings.NewReader(`{"id": 10, "name": "foo"}`)),
			},
			matcher: Matcher{
				BodyJSON: map[string]interface{}{"name": "foo"},
			},
		},
	}

	for _, d := range data {
		d := d
		b.Run(d.title, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = matchBodyJSON(d.req, d.matcher)
			}
		})
	}
}
