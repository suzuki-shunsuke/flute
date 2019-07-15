# fagott

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/suzuki-shunsuke/fagott/fagott)
[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/fagott/status.svg)](https://cloud.drone.io/suzuki-shunsuke/fagott)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/fagott/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/fagott)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/fagott)](https://goreportcard.com/report/github.com/suzuki-shunsuke/fagott)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/fagott.svg)](https://github.com/suzuki-shunsuke/fagott)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/fagott.svg)](https://github.com/suzuki-shunsuke/fagott/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/fagott/master/LICENSE)

Golang HTTP Client testing framework

## Presentation

https://speakerdeck.com/szksh/introduction-of-fagott-golang-http-client-testing-framework

## Overview

`fagott` is the Golang HTTP client testing framework.
The goal is

* Test request parameters such as the request path, headers and body
* Mock the HTTP server

`*fagott.Transport` implements [http.RoundTripper](https://golang.org/pkg/net/http/#RoundTripper).

`fagott` uses [testify](https://github.com/stretchr/testify)'s assert internally.
You can test the http request parameters with assert.

For example, the following test failure message means the request header is unexpected value.

```console
=== RUN   TestClient_CreateUser
--- FAIL: TestClient_CreateUser (0.00s)
    tester.go:168:
                Error Trace:    tester.go:168
                                                        tester.go:32
                                                        transport.go:25
                                                        client.go:250
                                                        client.go:174
                                                        client.go:641
                                                        client.go:509
                                                        create_user.go:45
                                                        create_user_test.go:56
                Error:          Not equal:
                                expected: []string{"token XXXXX"}
                                actual  : []string{"token "}

                                Diff:
                                --- Expected
                                +++ Actual
                                @@ -1,3 +1,3 @@
                                 ([]string) (len=1) {
                                - (string) (len=11) "token XXXXX"
                                + (string) (len=6) "token "
                                 }
                Test:           TestClient_CreateUser
                Messages:       the request header "Authorization" should match
                                service: http://example.com
                                request name: create a user
```

## Example

Please see [examples](examples).

## License

[MIT](LICENSE)
