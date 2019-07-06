package fagott

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRequest(t *testing.T, req *http.Request, service *Service, route *Route) {
	tester := route.Tester
	if tester.Path != "" {
		testPath(t, req, service, route)
	}
	if tester.Method != "" {
		testMethod(t, req, service, route)
	}
	if tester.BodyString != "" {
		testBodyString(t, req, service, route)
	}
	if tester.BodyJSON != nil {
		testBodyJSON(t, req, service, route)
	}
	if tester.Header != nil {
		testHeader(t, req, service, route)
	}
	if tester.Test != nil {
		tester.Test(t, req, service, route)
	}
}

func makeMsg(msg, srv, reqName string) string {
	return fmt.Sprintf(`%s
service: %s
request name: %s`, msg, srv, reqName)
}

func testBodyString(
	t *testing.T, req *http.Request, service *Service, route *Route,
) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	if req.Body == nil {
		assert.Equal(
			t, tester.BodyString, nil,
			makeMsg("request body should match", srv, reqName))
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to read the request body: %v", err),
				srv, reqName))
		return
	}
	assert.Equalf(
		t, tester.BodyString, string(b),
		makeMsg("request body should match", srv, reqName))
}

func testPath(t *testing.T, req *http.Request, service *Service, route *Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	assert.Equal(
		t, tester.Path, req.URL.Path,
		makeMsg("request path should match", srv, reqName))
}

func testMethod(t *testing.T, req *http.Request, service *Service, route *Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	assert.Equal(
		t, tester.Method, req.Method,
		makeMsg("request method should match", srv, reqName))
}

func testBodyJSON(t *testing.T, req *http.Request, service *Service, route *Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	if req.Body == nil {
		assert.Equal(
			t, tester.BodyJSON, nil,
			makeMsg("request body should match", srv, reqName))
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to read the request body: %v", err), srv, reqName))
		return
	}
	c, err := json.Marshal(tester.BodyJSON)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to parse tester.bodyJSON as JSON: %v", err),
				srv, reqName))
		return
	}
	assert.JSONEqf(
		t, string(b), string(c),
		makeMsg("request body should match", srv, reqName))
}

func testHeader(t *testing.T, req *http.Request, service *Service, route *Route) {
	reqName := route.Name
	srv := service.Endpoint

	for k, v := range route.Tester.Header {
		if v == nil {
			if _, ok := req.Header[k]; !ok {
				assert.Fail(
					t, makeMsg(
						`the following request header is required: `+k, srv, reqName))
				return
			}
		} else {
			a, ok := req.Header[k]
			if !ok {
				assert.Fail(
					t, makeMsg(
						"the following request header is required: "+k, srv, reqName))
				return
			}
			assert.Equal(
				t, v, a,
				makeMsg(fmt.Sprintf(`the request header "%s" should match`, k), srv, reqName))
		}
	}
}
