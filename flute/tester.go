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

func testRequest(t *testing.T, req *http.Request, service Service, route Route) {
	tester := route.Tester
	if tester == nil {
		// SPEC if the tester is nil, do nothing.
		return
	}
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
	if tester.BodyJSONString != "" {
		testBodyJSONString(t, req, service, route)
	}
	if tester.PartOfHeader != nil {
		testPartOfHeader(t, req, service, route)
	}
	if tester.Header != nil {
		assert.Equal(
			t, tester.Header, req.Header,
			makeMsg("request header should match", service.Endpoint, route.Name))
	}
	if tester.PartOfQuery != nil {
		testPartOfQuery(t, req, service, route)
	}
	if tester.Query != nil {
		assert.Equal(
			t, tester.Query, req.URL.Query(),
			makeMsg("request query parameter should match", service.Endpoint, route.Name))
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

func testBodyString(t *testing.T, req *http.Request, service Service, route Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	if req.Body == nil {
		assert.Equal(
			t, tester.BodyString, "",
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

func testPath(t *testing.T, req *http.Request, service Service, route Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	assert.Equal(
		t, tester.Path, req.URL.Path,
		makeMsg("request path should match", srv, reqName))
}

func testMethod(t *testing.T, req *http.Request, service Service, route Route) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	assert.Equal(
		t, strings.ToUpper(tester.Method), strings.ToUpper(req.Method),
		makeMsg("request method should match", srv, reqName))
}

func testBodyJSON(t *testing.T, req *http.Request, service Service, route Route) {
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
		t, string(c), string(b),
		makeMsg("request body should match", srv, reqName))
}

func testBodyJSONString(
	t *testing.T, req *http.Request, service Service, route Route,
) {
	reqName := route.Name
	srv := service.Endpoint
	tester := route.Tester

	if req.Body == nil {
		assert.Equal(
			t, tester.BodyString, "",
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
	assert.JSONEqf(
		t, tester.BodyJSONString, string(b),
		makeMsg("request body should match", srv, reqName))
}

func testPartOfHeader(t *testing.T, req *http.Request, service Service, route Route) {
	reqName := route.Name
	srv := service.Endpoint

	for k, v := range route.Tester.PartOfHeader {
		a, ok := req.Header[k]
		if !ok {
			assert.Fail(
				t, makeMsg(
					"the following request header is required: "+k, srv, reqName))
			return
		}
		if v != nil {
			assert.Equal(
				t, v, a,
				makeMsg(fmt.Sprintf(`the request header "%s" should match`, k), srv, reqName))
		}
	}
}

func testPartOfQuery(t *testing.T, req *http.Request, service Service, route Route) {
	reqName := route.Name
	srv := service.Endpoint

	query := req.URL.Query()
	for k, v := range route.Tester.PartOfQuery {
		a, ok := query[k]
		if !ok {
			assert.Fail(
				t, makeMsg(
					"the following request query is required: "+k, srv, reqName))
			return
		}
		if v != nil {
			assert.Equal(
				t, v, a,
				makeMsg(fmt.Sprintf(`the request query "%s" should match`, k), srv, reqName))
		}
	}
}
