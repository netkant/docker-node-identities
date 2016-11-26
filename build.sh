#!/bin/bash
BIN_NAME="cuddy"
GIT_REPO="github.com/urlund/cuddy"

if [ ! -z $1 ]; then
    # check if go is installed
    if [ $(which go > /dev/null 2>&1; echo $?) -ne 0 ]; then
        echo "you must have go installed to use this script"
        exit 1;
    fi

    # now build
    for GOOS in darwin linux; do
        for GOARCH in 386 amd64; do
            BIN_PATH="bin/$1/$BIN_NAME-$GOOS-$GOARCH"
            echo "building $BIN_PATH"
            go build -o $BIN_PATH
        done
    done
else
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
    TIMESTAMP=$(date +"%s")
    docker run --rm -it -v "$GOPATH":/work -e "GOPATH=/work" -w /work/src/$GIT_REPO golang:latest ./build.sh $TIMESTAMP
fi
