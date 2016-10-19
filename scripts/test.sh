#!/bin/bash

# Check if the Go runtime is installed.
if ! [ -x "$(which go)" ]; then
  echo "Go is not installed." >&2
  exit 1
fi

# Check if go-bindata is installed.
if ! [ -x "$(which golint)" ]; then
  echo "golint is not installed, installing..."
  go get -u github.com/golang/lint/golint
  echo "Installed."
fi

packages="go list ./... | grep -v '/vendor/' -v 'corald/fixtures/json'"

folders=(cmd internal)

for folder in "${folders[@]}"
do
  pushd $folder
    for package in $($package)
    do
      golint -set_exit_status $package
    done

    go vet ./...
    go test ./...
  popd
done
