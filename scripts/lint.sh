#!/bin/bash

pkgs=$(go list ./... | grep -v vendor)

echo "* Linting code for errors"

# Vet first since golangci-lint returns unclear errors if packages don't build
go vet ${pkgs}

# Exit if vet throws an error
vet_status="$?"

if [ "$vet_status" -ne 0 ]
then
	exit "$vet_status"
fi

golangci-lint run

echo "+ Your code is beautiful!"

exit 0
