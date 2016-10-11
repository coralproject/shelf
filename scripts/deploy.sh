#! /bin/bash

#
# README
#
# Running deploy.sh requires the following details available in the environment:
#
# - BUILD_NUMNER
# - GIT_COMMIT
# - AWS_DEFAULT_REGION
# - REPO_URL
#

# Based on an amazing article: http://dev.mikamai.com/post/144501132499/continuous-delivery-with-travis-and-ecs

set -v

# This is needed to login on AWS and push the image on ECR
# Change it accordingly to your docker repo
pip install --user awscli
export PATH=$PATH:$HOME/.local/bin
eval $(aws ecr get-login --region $AWS_DEFAULT_REGION)

CURRENT_DATE=$(date +%Y-%m-%d)
APPS=(askd corald sponged xeniad)

# For each app to be built, we will build the binary, and build the image,
# and then push it off to the registry.
for APP in "${APPS[@]}"
do
  IMAGE_NAME=coralproject/$APP
  REMOTE_IMAGE_URL="$REPO_URL/$IMAGE_NAME"

  echo "Building the $APP binary"
  GOOS=linux CGO_ENABLED=0 go build -a -ldflags "-X main.GitVersion=$BUILD_NUMNER -X main.GitRevision=$GIT_COMMIT -X main.BuildDate=$CURRENT_DATE" -o cmd/$APP/$APP github.com/coralproject/shelf/cmd/$APP

  echo "Build the $APP docker container"
  docker build -t $IMAGE_NAME cmd/$APP/

  echo "Pushing $IMAGE_NAME:latest"
  docker tag $IMAGE_NAME:latest "$REMOTE_IMAGE_URL:latest"
  docker push "$REMOTE_IMAGE_URL:latest"
  echo "Pushed $IMAGE_NAME:latest"
done
