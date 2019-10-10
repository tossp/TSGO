# TSGO

```bash
go fmt ./... && go build -v -ldflags "-s -w" && sites
go fmt ./... && go build -gcflags "all=-N -l" &&  dlv --listen=:2345 --headless=true --api-version=2 exec ./sites.exe
go get -u all
go mod edit -replace github.com/tossp/tsgo=../tsgo
```

[参考来源](https://github.com/golang-standards/project-layout)