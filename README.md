# bobber

Bobber is a tool for mocking out APIs your application integrates with. It replaces real APIs in a way that lets your application keep working.
- Records all incoming request details so it's easy to see exactly what you're sending
- Avoid spamming real APIs
- Test without real data - sometimes setting up real valid data can be a pain
- Control over responses - don't need to set up the perfect request

At any point Bobber can be removed and replaced with a real API. Bobber can also be setup to forward requests to APIs at the switch of a button to make testing and debugging a breeze.

## Install

### go install

This is preferred since it will automatically add `bobber` to your path

First [install Go](https://go.dev/doc/install) if you haven't already. Then run

```shell
go install github.com/tyerdimon/bobber
```

### Executables

Download from [GitHub]()

## Setup

### Namespaces

First create a namespace. If it helps you can think of each namespace as a different API to be mocked out. This helps avoid conflicts with common endpoints like `/users` or `/token`.

Each name space will get a slug, for example `really-cool-api`. A namespace with that slug will expect requests at `/requests/really-cool-api/`.

## Testing

This project includes both Javascript and Go tests.

### Javascript

```
npm ci
npm test
```