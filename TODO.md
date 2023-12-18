## TODO
- remove sqlx
- setup air or a auto reloader
- setup debugger. delve i think?
- logging
- error handling
- default sorting on queries
- TESTS!!!
- finish CRUD for everything

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