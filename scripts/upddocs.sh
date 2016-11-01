#!/bin/bash

# Check if the Go runtime is installed.
if ! [ -x "$(which go)" ]; then
  echo "Go is not installed." >&2
  exit 1
fi

# Check if go-bindata is installed.
if ! [ -x "$(which godoc2md)" ]; then
  echo "godoc2md is not installed, installing..."
  go get -u github.com/davecheney/godoc2md
  echo "Installed."
fi

# list all the packges, trim out the vendor directory and any main packages,
# then strip off the package name
PACKAGES="$(go list -f "{{.Name}}:{{.ImportPath}}" ./... | grep -v -E "main:|vendor/" | cut -d ":" -f 2)"

# loop over all packages generating all their documentation
for PACKAGE in $PACKAGES
do

  echo "godoc2md $PACKAGE > $GOPATH/src/$PACKAGE/README.md"

  godoc2md $PACKAGE -links=false > $GOPATH/src/$PACKAGE/README.md

done
