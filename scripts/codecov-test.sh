#!/usr/bin/env bash
# https://github.com/codecov/example-go#caveat-multiple-files

echo "" > coverage.txt

# the example package is tested but the coverage is ignored.
go test -race -covermode=atomic ./examples || exit 1
for d in $(go list ./... | grep -v vendor | grep -v examples); do
  go test -race -coverprofile=profile.out -covermode=atomic $d || exit 1
  if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
  fi
done
