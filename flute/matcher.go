package flute

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/suzuki-shunsuke/go-dataeq/v2/dataeq"
)

// isMatchService returns whether the request matches with the service.
// isMatchService checks the request URL.Scheme and URL.Host are equal to the service endpoint.
func isMatchService(req *http.Request, service Service) bool {
	return req.URL.Scheme+"://"+req.URL.Host == service.Endpoint
}

type matchFunc func(req *http.Request, matcher Matcher) (bool, error)

func matchPath(req *http.Request, matcher Matcher) (bool, error) {
	return matcher.Path == "" || matcher.Path == req.URL.Path, nil
}

func matchMethod(req *http.Request, matcher Matcher) (bool, error) {
	return matcher.Method == "" || strings.EqualFold(matcher.Method, req.Method), nil
}

func matchHeader(req *http.Request, matcher Matcher) (bool, error) {
	return matcher.Header == nil || reflect.DeepEqual(matcher.Header, req.Header), nil
}

func matchQuery(req *http.Request, matcher Matcher) (bool, error) {
	return matcher.Query == nil || reflect.DeepEqual(matcher.Query, req.URL.Query()), nil
}

var matchFuncs = [...]matchFunc{ //nolint:gochecknoglobals
	matchPath, matchMethod, matchBodyString, matchBodyJSON, matchBodyJSONString,
	matchPartOfHeader, matchHeader, matchPartOfQuery, matchQuery,
}

// isMatch returns whether the request matches with the matcher.
// If the matcher has multiple conditions, IsMatch returns true if the request meets all conditions.
func isMatch(req *http.Request, matcher Matcher) (bool, error) {
	for _, match := range matchFuncs {
		if f, err := match(req, matcher); err != nil || !f {
			return f, err
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

func matchPartOfHeader(req *http.Request, matcher Matcher) (bool, error) {
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

func matchPartOfQuery(req *http.Request, matcher Matcher) (bool, error) {
	if matcher.PartOfQuery == nil {
		return true, nil
	}
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

func matchBodyString(req *http.Request, matcher Matcher) (bool, error) {
	if matcher.BodyString == "" {
		return true, nil
	}
	if req.Body == nil {
		return false, nil
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	return matcher.BodyString == string(b), nil
}

func matchBodyJSONString(req *http.Request, matcher Matcher) (bool, error) {
	if matcher.BodyJSONString == "" {
		return true, nil
	}
	if req.Body == nil {
		return false, nil
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	return dataeq.JSON.Equal(b, []byte(matcher.BodyJSONString))
}

func matchBodyJSON(req *http.Request, matcher Matcher) (bool, error) {
	if matcher.BodyJSON == nil {
		return true, nil
	}
	if req.Body == nil {
		return false, nil
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read the request body: %w", err)
	}
	return dataeq.JSON.Equal(b, matcher.BodyJSON)
}
