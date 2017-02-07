#!/bin/bash
VERSION=`cat VERSION`
BIN_NAME="docker-node-identities"
BIN_PATH="bin/$BIN_NAME-$VERSION"
GIT_REPO="github.com/urlund/$BIN_NAME"

# check if docker is installed
if [ $(which docker > /dev/null 2>&1; echo $?) -ne 0 ]; then
    echo "you must have docker installed to use this script"
    exit 1;
fi

# check if $GOPATH isset
if [ -z $GOPATH ]; then
    echo "you must set \$GOPATH to use this script"
    exit 1;
fi

# run docker with build cmd
docker run --rm -it -v "$GOPATH":/work -e "GOPATH=/work" -w /work/src/$GIT_REPO golang:latest go build -o $BIN_PATH
