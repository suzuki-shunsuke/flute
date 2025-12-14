# flute

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![DeepWiki](https://img.shields.io/badge/DeepWiki-suzuki--shunsuke%2Fflute-blue.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAyCAYAAAAnWDnqAAAAAXNSR0IArs4c6QAAA05JREFUaEPtmUtyEzEQhtWTQyQLHNak2AB7ZnyXZMEjXMGeK/AIi+QuHrMnbChYY7MIh8g01fJoopFb0uhhEqqcbWTp06/uv1saEDv4O3n3dV60RfP947Mm9/SQc0ICFQgzfc4CYZoTPAswgSJCCUJUnAAoRHOAUOcATwbmVLWdGoH//PB8mnKqScAhsD0kYP3j/Yt5LPQe2KvcXmGvRHcDnpxfL2zOYJ1mFwrryWTz0advv1Ut4CJgf5uhDuDj5eUcAUoahrdY/56ebRWeraTjMt/00Sh3UDtjgHtQNHwcRGOC98BJEAEymycmYcWwOprTgcB6VZ5JK5TAJ+fXGLBm3FDAmn6oPPjR4rKCAoJCal2eAiQp2x0vxTPB3ALO2CRkwmDy5WohzBDwSEFKRwPbknEggCPB/imwrycgxX2NzoMCHhPkDwqYMr9tRcP5qNrMZHkVnOjRMWwLCcr8ohBVb1OMjxLwGCvjTikrsBOiA6fNyCrm8V1rP93iVPpwaE+gO0SsWmPiXB+jikdf6SizrT5qKasx5j8ABbHpFTx+vFXp9EnYQmLx02h1QTTrl6eDqxLnGjporxl3NL3agEvXdT0WmEost648sQOYAeJS9Q7bfUVoMGnjo4AZdUMQku50McDcMWcBPvr0SzbTAFDfvJqwLzgxwATnCgnp4wDl6Aa+Ax283gghmj+vj7feE2KBBRMW3FzOpLOADl0Isb5587h/U4gGvkt5v60Z1VLG8BhYjbzRwyQZemwAd6cCR5/XFWLYZRIMpX39AR0tjaGGiGzLVyhse5C9RKC6ai42ppWPKiBagOvaYk8lO7DajerabOZP46Lby5wKjw1HCRx7p9sVMOWGzb/vA1hwiWc6jm3MvQDTogQkiqIhJV0nBQBTU+3okKCFDy9WwferkHjtxib7t3xIUQtHxnIwtx4mpg26/HfwVNVDb4oI9RHmx5WGelRVlrtiw43zboCLaxv46AZeB3IlTkwouebTr1y2NjSpHz68WNFjHvupy3q8TFn3Hos2IAk4Ju5dCo8B3wP7VPr/FGaKiG+T+v+TQqIrOqMTL1VdWV1DdmcbO8KXBz6esmYWYKPwDL5b5FA1a0hwapHiom0r/cKaoqr+27/XcrS5UwSMbQAAAABJRU5ErkJggg==)](https://deepwiki.com/suzuki-shunsuke/flute) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/suzuki-shunsuke/flute/flute)
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
