#!/bin/bash

set -e

GO_FLAGS=${GO_FLAGS:-"-tags netgo"}
GO_CMD=${GO_CMD:-"build"}
BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}
VERBOSE=${VERBOSE:-}

app_name="leaf"
repo_path="github.com/vrongmeal/leaf"

# Get branch revision and  version
version="1.0.1"
revision=$(git rev-parse --short HEAD 2> /dev/null || echo 'unknown')
branch=$(git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown')

# Extract the go version
go_version=$(go version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')

# go 1.4 requires ldflags format to be "-X key value", not "-X key=value"
# ldseparator here is for cross compatibility between go versions
ldseparator="="
if [ "${go_version:0:3}" = "1.4" ]; then
	ldseparator=" "
fi

ldflags="
  -X ${repo_path}/version.AppName${ldseparator}${app_name}
  -X ${repo_path}/version.Version${ldseparator}${version}
  -X ${repo_path}/version.Revision${ldseparator}${revision}
  -X ${repo_path}/version.Branch${ldseparator}${branch}
  -X ${repo_path}/version.BuildUser${ldseparator}${BUILD_USER}
  -X ${repo_path}/version.BuildDate${ldseparator}${BUILD_DATE}
  -X ${repo_path}/version.GoVersion${ldseparator}${go_version}"

echo "* Building into build/${app_name}"

if [ -n "$VERBOSE" ]; then
  echo "* Building with -ldflags $ldflags"
fi

GOBIN=$PWD go "${GO_CMD}" -o "${PWD}/build/${app_name}" ${GO_FLAGS} -ldflags "${ldflags}" "${repo_path}"

echo "+ Build complete!"

exit 0
