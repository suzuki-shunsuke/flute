package fagott

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/pkg/errors"
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
	if matcher.Path != "" {
		if matcher.Path != req.URL.Path {
			return false, nil
		}
	}
	if matcher.Method != "" {
		if strings.ToUpper(matcher.Method) != strings.ToUpper(req.Method) {
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
	if matcher.Header != nil {
		for k, v := range matcher.Header {
			if v == nil {
				if _, ok := req.Header[k]; !ok {
					return false, nil
				}
			} else {
				a, ok := req.Header[k]
				if !ok {
					return false, nil
				}
				if !reflect.DeepEqual(a, v) {
					return false, nil
				}
			}
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

func isMatchBodyString(req *http.Request, matcher *Matcher) (bool, error) {
	if req.Body == nil {
		return false, nil
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to read the request body")
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
		return false, errors.Wrap(err, "failed to read the request body")
	}
	return jsoneq.Equal(b, []byte(matcher.BodyJSONString))
}

func isMatchBodyJSON(req *http.Request, matcher *Matcher) (bool, error) {
	if req.Body == nil {
		return false, nil
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, errors.Wrap(err, "failed to read the request body")
	}
	return jsoneq.Equal(b, matcher.BodyJSON)
}
