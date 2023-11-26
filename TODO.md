## TODO

### FIX

- clean up CSS
- finish go tests
- set up .env file

### NEW

 - add namespaces
   - part of path following `/requests/` before the `next` will map to a namespace
 - add endpoints
   - path to match again incoming requests
   - response
 - import / export namespaces and endpoints to JSON
 - dynamic responses
   - dynamic based off path parts, split path on `/`
   - dynamic based off request body, `{level1.level2.data}` grabs value from request
 - dynamic response headers
 - better support for form data
