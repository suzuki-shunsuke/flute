# flute

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/suzuki-shunsuke/flute/flute)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/flute)](https://goreportcard.com/report/github.com/suzuki-shunsuke/flute)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/flute/main/LICENSE)

Golang HTTP client testing framework

## Presentation

https://speakerdeck.com/szksh/flute-golang-http-client-testing-framework

## Overview

`flute` is the Golang HTTP client testing framework.
The goal is

* Test request parameters such as the request path, headers and body
* Mock the HTTP server

`flute.Transport` implements [http.RoundTripper](https://golang.org/pkg/net/http/#RoundTripper).

`flute` uses [testify](https://github.com/stretchr/testify)'s assert internally.
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

Please see [the example](https://github.com/suzuki-shunsuke/flute/blob/v0.6.0/examples/create_user_test.go#L21-L48).

## Example

Please see [examples](examples).

## License

[MIT](LICENSE)
