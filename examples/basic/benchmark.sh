#!/bin/bash

go build -o char ../../main.go

function devcharls() {
  CHAR_DEV_MODE=true ./char ls 2&> /dev/null
}

function charls() {
  ./char ls 2&> /dev/null
}

echo "scratch"
time devcharls
echo ""

echo "ready"
time charls
echo ""
