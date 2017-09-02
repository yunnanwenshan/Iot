#!/usr/bin/env bash

set -e

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
OLDGOBIN="$GOBIN"
export GOPATH="$CURDIR"
export GOBIN="$CURDIR/bin/"
echo 'GOPATH:' $GOPATH
echo 'GOBIN:' $GOBIN

if [ -e ./app.pid ]; then
	#PID=$(cat ./app.pid)
	#echo "PID:" $PID
	PID=$(ps aux | grep ./bin/app | grep -v grep | awk '{print $2}')
	echo $PID
	if [ -n "$PID" ];then
		kill -1 $PID
    fi
else
    nohup ./bin/app
fi

export GOPATH="$OLDGOPATH"
export GOBIN="$OLDGOBIN"

echo 'run finished'

