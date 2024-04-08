#!/bin/sh

version="v1.54.1"
gobin=${GOPATH}/bin
binpath=${gobin}/golangci-lint-${version}

if [ ! -f ${binpath} ]; then
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@${version}
  mv ${gobin}/golangci-lint ${binpath}
fi

${binpath} run -v -c .golangci.yml ./...
