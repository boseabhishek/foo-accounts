# foo-accounts

gRPC + HTTP Proxy + API Docs - all done via Go.

## generate grpc stuff:

```shell
    buf generate
```

## run the gRPC/HTTP service:

```shell
    go mod tidy
    go run cmd/server/main.go
```