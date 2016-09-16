#!/bin/bash

# Check if the Go runtime is installed.
if ! [ -x "$(which go)" ]; then
  echo "Go is not installed." >&2
fi

# Check if go-bindata is installed.
if ! [ -x "$(which go-bindata)" ]; then
  echo "go-bindata is not installed, installing..."
  go get -u github.com/jteeuwen/go-bindata/...
  echo "Installed."
fi

# Run go generate
echo "Generating asset fixtures..."
go generate github.com/coralproject/shelf/cmd/corald/...
echo "Generated."
