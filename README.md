# HTTP calls from NodeJS to Golang

## Contents

`src/index.ts`

ExpressJS server and handler for incoming calls.

`napi/gohandler/src/addon.cc`

NodeJS N-API interop package.

`cmd/gohandler/main.go`

HTTP Handler CGO export.

`internal/handlers/handlers.go`

Go http ServeHTTP handler.
