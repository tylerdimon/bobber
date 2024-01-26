## TODO

- support for headers
- tests
- endpoint edit, delete
- setup auto reload. air? modd?
- setup debugger. delve?
- import / export namespaces and endpoints to JSON

- think about routing in general
    - listener /{apiname}/{endpoint}
    - config /config/{apiname}/{endpoint-slug}
    - ui / -> /ui/{apiname}/{endpoint-slug} (because listener should never need / url, we should redirect to UI for better UX)
- TLDR the pages at /requests and / should switch.

 - dynamic responses
   - dynamic based off path parts, split path on `/`
   - dynamic based off request body, `{level1.level2.data}` grabs value from request

 - dynamic response headers
 - better support for form data