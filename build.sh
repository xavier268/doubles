#!/bin/bash

VERSION="v1.1
"
git commit -a -m "updating binaries"

GIT_HASH=$(git log --format=%H -1 HEAD)
GIT_DATE=$(git log --format=%cD -1 HEAD)
VERSION="$VERSION - Build $GIT_HASH ($GIT_DATE)"

echo "Last commit full hash : $GIT_HASH"
echo "Last commit date : $GIT_DATE"
echo $VERSION


echo "buiding linux binary"
go build -ldflags "-X 'main.version=$VERSION'"  .

echo "Buiding windows binary"
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.version=$VERSION'"  .
