# go-grpc-service

workload for grpc tests

## server

server grpc workload 

## client

client grpc workload for test

    k port-forward svc/svc-go-grpc-service-server -n test-b 50051:50051

## client-http-redirect

webserver receive rest call and redirect to server grpc

## Compile

    protoc -I proto proto/service.proto --go_out=plugins=grpc:proto
