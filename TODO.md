### TODO
- HTML for requests is in two places, websocket message should just send HTML 
rendered by Go to keep it in same place?
- finish go tests
- set up .env and move things to it


 - add namespaces
   - part of path following `/requests/` before the `next` will map to a namespace
 - add endpoints
   - path to match again incoming requests
   - response
 - import / export namespaces and endpoints to JSON
 - dynamic responses
   - dynamic based off path parts, split path on `/`
   - dynamic based off requst body, `{level1.level2.data}` grabs value from request
 - forwarding - a checkbox that sends requests to another API when enabled
   - each namespace needs a provided URL
   - without modifying but still recording request forward on to URL
   - what about authentication?
 - dynamic response headers
 - better support for form data
 - websockets? other protocols?
