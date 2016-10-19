#! /bin/bash

#
# README
#
# Running deploy.sh requires the following details available in the environment:
#
# - BUILD_NUMNER
# - GIT_COMMIT
#

# Based on an amazing article: http://dev.mikamai.com/post/144501132499/continuous-delivery-with-travis-and-ecs

set -v

CURRENT_DATE=$(date +%Y-%m-%d)
APPS=(askd corald sponged xeniad)

# For each app to be built, we will build the binary, and build the image,
# and then push it off to the registry.
for APP in "${APPS[@]}"
do
  IMAGE_NAME=coralproject/$APP

  echo "Building the $APP binary"
  GOOS=linux CGO_ENABLED=0 go build -a -ldflags "-X main.GitVersion=$BUILD_NUMNER -X main.GitRevision=$GIT_COMMIT -X main.BuildDate=$CURRENT_DATE" -o cmd/$APP/$APP github.com/coralproject/shelf/cmd/$APP

  echo "Build the $APP docker container"
  docker build -t $IMAGE_NAME cmd/$APP/

  docker tag $IMAGE_NAME:latest

  echo "Pushing $IMAGE_NAME:latest"
  docker push $IMAGE_NAME:latest
  echo "Pushed $IMAGE_NAME:latest"
done
