#!/bin/sh

echo "Starting proto to struct..."

protoc -I=. -I=$GOPATH/src/ --go_out=.. --gorm_out=engine=postgres:.. *.proto

echo "Completed proto to struct..."
