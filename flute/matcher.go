package flute

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/suzuki-shunsuke/go-jsoneq/jsoneq"
)

// isMatchService returns whether the request matches with the service.
// isMatchService checks the request URL.Scheme and URL.Host are equal to the service endpoint.
func isMatchService(req *http.Request, service *Service) bool {
	return req.URL.Scheme+"://"+req.URL.Host == service.Endpoint
}

// isMatch returns whether the request matches with the matcher.
// If the matcher has multiple conditions, IsMatch returns true if the request meets all conditions.
func isMatch(req *http.Request, matcher *Matcher) (bool, error) {
	if matcher == nil {
		// SPEC if the matcher is nil, the route matches the request.
		return true, nil
	}
	if matcher.Path != "" {
		if matcher.Path != req.URL.Path {
			return false, nil
		}
	}
	if matcher.Method != "" {
		if !strings.EqualFold(matcher.Method, req.Method) {
			return false, nil
		}
	}
	if matcher.BodyString != "" {
		f, err := isMatchBodyString(req, matcher)
		if err != nil || !f {
			return f, err
		}
	}
	if matcher.BodyJSON != nil {
		f, err := isMatchBodyJSON(req, matcher)
		if err != nil || !f {
			return f, err
		}
	}
	if matcher.BodyJSONString != "" {
		f, err := isMatchBodyJSONString(req, matcher)
		if err != nil || !f {
			return f, err
		}
	}
	if matcher.PartOfHeader != nil {
		f, err := isMatchPartOfHeader(req, matcher)
		if err != nil || !f {
			return f, err
		}
	}
	if matcher.Header != nil {
		if !reflect.DeepEqual(matcher.Header, req.Header) {
			return false, nil
		}
	}
	if matcher.PartOfQuery != nil {
		f, err := isMatchPartOfQuery(req, matcher)
		if err != nil || !f {
			return f, err
		}
	}
	if matcher.Query != nil {
		if !reflect.DeepEqual(matcher.Query, req.URL.Query()) {
			return false, nil
		}
	}
	if matcher.Match != nil {
		f, err := matcher.Match(req)
		if err != nil || !f {
			return f, err
		}
	}
	return true, nil
}

func isMatchPartOfHeader(req *http.Request, matcher *Matcher) (bool, error) {
	for k, v := range matcher.PartOfHeader {
		a, ok := req.Header[k]
		if !ok {
			return false, nil
		}
		if v != nil {
			if !reflect.DeepEqual(a, v) {
				return false, nil
			}
		}
	}
	return true, nil
}

func isMatchPartOfQuery(req *http.Request, matcher *Matcher) (bool, error) {
	query := req.URL.Query()
	for k, v := range matcher.PartOfQuery {
		a, ok := query[k]
		if !ok {
			return false, nil
		}
		if v != nil {
			if !reflect.DeepEqual(a, v) {
				return false, nil
			}
		}
	}
	return true, nil
}

func isMatchBodyString(req *http.Request, matcher *Matcher) (bool, error) {
	if req.Body == nil {
		return false, nil
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	if matcher.BodyString != string(b) {
		return false, nil
	}
	return true, nil
}

func isMatchBodyJSONString(req *http.Request, matcher *Matcher) (bool, error) {
	if req.Body == nil {
		return false, nil
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	return jsoneq.Equal(b, []byte(matcher.BodyJSONString))
}

func isMatchBodyJSON(req *http.Request, matcher *Matcher) (bool, error) {
	if req.Body == nil {
		return false, nil
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	return jsoneq.Equal(b, matcher.BodyJSON)
}
