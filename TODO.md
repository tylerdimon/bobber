## TODO

- namespaces
  - logging
  - error handling
  - default sorting on queries
  - remove sqlx usage
  - add tests
    - http
    - db
- requests
- endpoints

- setup auto reload. air? modd?
- setup debugger. delve?

- think about routing in general
    - listener /{apiname}/{endpoint}
    - config /config/{apiname}/{endpoint-slug}
    - ui / -> /ui/{apiname}/{endpoint-slug} (because listener should never need / url, we should redirect to UI for better UX)

 - import / export namespaces and endpoints to JSON

 - dynamic responses
   - dynamic based off path parts, split path on `/`
   - dynamic based off request body, `{level1.level2.data}` grabs value from request

 - dynamic response headers
 - better support for form data