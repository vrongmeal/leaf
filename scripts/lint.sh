#!/bin/bash

set -e

pkgs=$(go list ./... | grep -v vendor)

echo "* Linting code for errors"

# Vet first since golangci-lint returns unclear errors if packages don't build
go vet ${pkgs}

golangci-lint run

echo "+ Your code is beautiful!"

exit 0
