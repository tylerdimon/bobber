# Using bobber

UI runs on `localhost:8000` by default

Listening for requests at `locahost:8000/requests/`

### Setup

The first step is creating a namespace.
Think of each namespace as a different API to be mocked out.
This helps avoid conflicts with common endpoints like `/users` or `/token`.

Each name space will get a slug, for example `really-cool-api`. A namespace with that slug will expect requests at `/requests/really-cool-api/`.
