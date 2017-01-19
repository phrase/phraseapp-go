#!/bin/bash
set -e -o pipefail

export GO_VERSION=1.7.4

docker run -v $(pwd):/go/src/github.com/phrase/phraseapp-go -w /go/src/github.com/phrase/phraseapp-go golang:${GO_VERSION} make all
