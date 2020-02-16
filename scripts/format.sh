#!/bin/bash

pkgs=$(go list ./... | grep -v vendor)

echo "* Formatting code"

go vet ${pkgs} > /dev/null 2>&1

vet_status="$?"

if [ "$vet_status" -ne 0 ]
then
	go fmt ${pkgs}
	echo "- Error while running golangci-lint formatters."
    echo "  Run 'make lint' to fix errors."
	exit 1
fi

# Run --fix with exit code 0 and don't display output of command
golangci-lint run --fix --issues-exit-code 0 > /dev/null 2>&1

echo "+ Format complete!"

exit 0
