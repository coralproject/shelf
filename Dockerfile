# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.7.1

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/coralproject/shelf

# Build & Install
RUN cd /go/src && go install github.com/coralproject/shelf/cmd/xeniad && go install github.com/coralproject/shelf/cmd/xenia

# Run the app
ENTRYPOINT ["/go/src/github.com/coralproject/shelf/entrypoint.sh"]

# Document that the service listens on port 8080.
EXPOSE 4000
