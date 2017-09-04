#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN
go build -o node -race -gcflags "-N -l"  src/iot/node/*.go

if [ ! -d ./bin/node ]; then
	mkdir -p bin/node
fi

if [ -e ./node ]; then
	mv node ./bin/node/
	cp src/iot/node/config.ini ./bin/node/config.ini
	cp src/iot/node/log.xml ./bin/node/log.xml
	echo "copy node to ./bin/node"
fi

go build -o router -race -gcflags "-N -l"  src/iot/router/*.go

if [ ! -d ./bin/router ]; then
	mkdir -p bin/router
fi

if [ -e ./router ]; then
    mv router ./bin/router/
	cp src/iot/router/config.ini ./bin/router/config.ini
	cp src/iot/router/log.xml ./bin/router/log.xml
	echo "copy router to ./bin/router"
fi

go build -o gateway-api -race -gcflags "-N -l"  src/iot/gateway/main.go

if [ -e ./gateway-api ]; then
     mv gateway-api ./bin/
     echo "copy gateway-api to ./bin/gateway-api"
fi

go build -o device -race -gcflags "-N -l"  src/iot/device/dev/dev.go
echo "device build finished"

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'build finished'

