## TODO

### FIX

- clean up CSS
- finish go tests
- set up .env file
- default sorting on queries
- fix null strings with sqlx

### NEW

- config
  - html
    - add new

- think about routing in general
    - listener /{apiname}/{endpoint}
    - config /config/{apiname}/{endpoint-slug}
    - ui / -> /ui/{apiname}/{endpoint-slug} (because listener should never need / url, we should redirect to UI for better UX)

- add endpoints
   - path to match again incoming requests
   - response
 - add support for multiple APIS / namespaces

 - import / export namespaces and endpoints to JSON

 - dynamic responses
   - dynamic based off path parts, split path on `/`
   - dynamic based off request body, `{level1.level2.data}` grabs value from request

 - dynamic response headers
 - better support for form data