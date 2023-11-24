# Bobber

Bobber is a tool for mocking out APIs your application integrates with using HTTP request. 
It replaces real APIs in a way that lets your application keep working while making development easier.

Swap a Bobber URL with an API URL to start sending requests. Bobber will record all requests and send a response.
Responses can be static or dynamic based on request data.

## Use

### Install

Work in progress, for now [building](#Build) is required

### Setup

The first step is creating a namespace.
Think of each namespace as a different API to be mocked out. 
This helps avoid conflicts with common endpoints like `/users` or `/token`.

Each name space will get a slug, for example `really-cool-api`. A namespace with that slug will expect requests at `/requests/really-cool-api/`.

## Develop

To build and run

```bash
go build ./cmd/bobber
./bobber
```

This project includes both Javascript and Go tests.

#### Javascript

```
npm ci
npm test
```

#### Go

```
go test ./...
```

### Formatting

Use `Gofmt` for formatting

```
go fmt ./...
```