#!/bin/bash
# golang version
protoc --go_out=proto --go-grpc_out=proto -I proto proto/grpc.proto

# gopherjs version
protoc -I proto/ proto/grpc.proto --gopherjs_out=plugins=grpc:proto/gopherjs

