## TODO

### FIX

- log the errors
- clean up error handling
- add tests


- default sorting on queries
- fix null strings with sqlx
- split config handler to better match patterns?
- some models should have created at / updated at instead of timestamp field
- make sure sqlite service is handling timestamps properly
- switch IDs over to UUID
- MustParse probably needs to be replaced

### NEW

- add namespace & endpoint logic to request handler
- update endpoint
- delete namespace
- delete endpoint
- add default response to namespace
- add database column mappings to structs
- some models should have created at / updated at instead of timestamp field

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