# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.5.3
ENV GO15VENDOREXPERIMENT="1"
# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/coralproject/xenia
# Build & Install
RUN cd /go/src && go install github.com/coralproject/xenia/cmd/xeniad
# Run the app
ENTRYPOINT /go/bin/xeniad
# Document that the service listens on port 8080.
EXPOSE 4000
