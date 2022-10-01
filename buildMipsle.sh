#!/bin/bash


GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build  -ldflags "-s -w" -o webat main.go
scp webat root@192.168.4.1:/tmp

scp -r html root@192.168.4.1:/tmp