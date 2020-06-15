package flute

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testFunc func(t *testing.T, req *http.Request, service Service, route Route)

var testFuncs = [...]testFunc{ //nolint:gochecknoglobals
	testPath, testMethod, testBodyString, testBodyJSON,
	testBodyJSONString, testPartOfHeader, testHeader, testPartOfQuery,
	testQuery,
}

func testHeader(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.Header == nil {
		return
	}
	assert.Equal(
		t, route.Tester.Header, req.Header,
		makeMsg("request header should match", service.Endpoint, route.Name))
}

func testQuery(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.Query == nil {
		return
	}
	assert.Equal(
		t, route.Tester.Query, req.URL.Query(),
		makeMsg("request query parameter should match", service.Endpoint, route.Name))
}

func testRequest(t *testing.T, req *http.Request, service Service, route Route) {
	for _, fn := range testFuncs {
		fn(t, req, service, route)
	}
	tester := route.Tester
	if tester.Test != nil {
		tester.Test(t, req, service, route)
	}
}

func makeMsg(msg, srv, reqName string) string {
	return fmt.Sprintf(`%s
service: %s
request name: %s`, msg, srv, reqName)
}

func testBodyString(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.BodyString == "" {
		return
	}

	if req.Body == nil {
		assert.Equal(
			t, route.Tester.BodyString, "",
			makeMsg("request body should match", service.Endpoint, route.Name))
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to read the request body: %v", err),
				service.Endpoint, route.Name))
		return
	}
	assert.Equalf(
		t, route.Tester.BodyString, string(b),
		makeMsg("request body should match", service.Endpoint, route.Name))
}

func testPath(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.Path == "" {
		return
	}
	assert.Equal(
		t, route.Tester.Path, req.URL.Path,
		makeMsg("request path should match", service.Endpoint, route.Name))
}

func testMethod(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.Method == "" {
		return
	}

	assert.Equal(
		t, strings.ToUpper(route.Tester.Method), strings.ToUpper(req.Method),
		makeMsg("request method should match", service.Endpoint, route.Name))
}

func testBodyJSON(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.BodyJSON == nil {
		return
	}

	if req.Body == nil {
		assert.Equal(
			t, route.Tester.BodyJSON, nil,
			makeMsg("request body should match", service.Endpoint, route.Name))
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to read the request body: %v", err), service.Endpoint, route.Name))
		return
	}
	c, err := json.Marshal(route.Tester.BodyJSON)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to parse route.Tester.bodyJSON as JSON: %v", err),
				service.Endpoint, route.Name))
		return
	}
	assert.JSONEqf(
		t, string(c), string(b),
		makeMsg("request body should match", service.Endpoint, route.Name))
}

func testBodyJSONString(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.BodyJSONString == "" {
		return
	}

	if req.Body == nil {
		assert.Equal(
			t, route.Tester.BodyString, "",
			makeMsg("request body should match", service.Endpoint, route.Name))
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		assert.Fail(
			t, makeMsg(
				fmt.Sprintf("failed to read the request body: %v", err),
				service.Endpoint, route.Name))
		return
	}
	assert.JSONEqf(
		t, route.Tester.BodyJSONString, string(b),
		makeMsg("request body should match", service.Endpoint, route.Name))
}

func testPartOfHeader(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.PartOfHeader == nil {
		return
	}

	for k, v := range route.Tester.PartOfHeader {
		a, ok := req.Header[k]
		if !ok {
			assert.Fail(
				t, makeMsg(
					"the following request header is required: "+k, service.Endpoint, route.Name))
			return
		}
		if v != nil {
			assert.Equal(
				t, v, a,
				makeMsg(fmt.Sprintf(`the request header "%s" should match`, k), service.Endpoint, route.Name))
		}
	}
}

func testPartOfQuery(t *testing.T, req *http.Request, service Service, route Route) {
	if route.Tester.PartOfQuery == nil {
		return
	}

	query := req.URL.Query()
	for k, v := range route.Tester.PartOfQuery {
		a, ok := query[k]
		if !ok {
			assert.Fail(
				t, makeMsg(
					"the following request query is required: "+k, service.Endpoint, route.Name))
			return
		}
		if v != nil {
			assert.Equal(
				t, v, a,
				makeMsg(fmt.Sprintf(`the request query "%s" should match`, k), service.Endpoint, route.Name))
		}
	}
}
