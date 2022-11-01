#!/bin/bash

go build -o char ../../main.go

CHAR_DEV_MODE=true ./char ls
