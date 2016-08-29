progname = viewworker
tag = $(REGISTRY_NAME)/$(progname)

all: compile docker push

compile:
		CGO_ENABLED=0 go build -o $(progname)

docker:
		sudo docker build -t $(registry)/$(progname) .

push:
		sudo docker push $(registry)/$(progname)
