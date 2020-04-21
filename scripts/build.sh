#!/bin/bash

set -e

echo "* Building as ./build/leaf"

test -d build || mkdir build
go build -o build/leaf cmd/leaf/main.go

echo "+ Build complete!"

exit 0
