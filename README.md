# Bobber

Bobber is a tool for mocking out APIs your application integrates with using HTTP request. 
It replaces real APIs in a way that lets your application keep working while making development easier.

Swap a Bobber URL with an API URL to start sending requests. Bobber will record all requests and send a response.
Responses can be static or dynamic based on request data.

The README focuses on development, for information on using bobber go to [doc.md](doc.md)

## Build & run

```bash
go build ./cmd/bobber
./bobber
```

UI runs on `localhost:8000` by default

Listening for requests at `locahost:8000/requests/`

## Testing

This project includes both Javascript and Go tests.

### Javascript

```
npm ci
npm test
```

### Go

Mocks

```
rm -rf mocks
mockery
```

Test

```
go install gotest.tools/gotestsum@latest
gotestsum --format dots
```

Coverage

```shell
 go test ./... -coverprofile=c.out
go tool cover -html=c.out
```

### Formatting

```
go fmt ./...
```