#!/bin/bash

echo "* Installing tools"

# golangci-lint
echo "* Installing golangci-lint"
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.22.2

echo "+ Install complete!"

exit 0
