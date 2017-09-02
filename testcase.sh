#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN
#go test api/service
go test -v -race -test.run Info api/http/controller
go test -race -test.v -test.run BikeVersion api/http/controller
go test -race -test.v -test.run BikeModel api/http/controller
go test -race -test.v -test.run CheckType api/http/controller
#go test -race -test.v -test.run Info -bench=. -benchtime=2s -cpuprofile=prof.cpu pi/http/controller
go test -race -test.v -test.run Info -bench=. -benchtime=2s -memprofile=prof.mem api/http/controller
#go test -test.v api/http/controller
export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'test finished'

