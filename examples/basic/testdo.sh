#!/bin/bash

set -e

go build -o char ../../main.go

echo "base"
CHAR_DEV_MODE=true ./char do -h

echo 
echo "--------"
echo "local_up"
CHAR_DEV_MODE=false ./char do local_up --fish something
