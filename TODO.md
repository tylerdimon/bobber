## TODO


------------------
- think about routing in general
    - listener /{apiname}/{endpoint}
    - config /config/{apiname}/{endpoint-slug}
    - ui / -> /ui/{apiname}/{endpoint-slug} (because listener should never need / url, we should redirect to UI for better UX)
- TLDR the pages at /requests and / should switch.

 - dynamic response headers

- package
- debug / prod mode for html/css
```go
staticFilesDebugMode := os.Getenv("DEBUG")
if staticFilesDebugMode == "True" {
    log.Println("Serving static files from directory...")
    s.router.PathPrefix("/static/assets").
        Handler(http.StripPrefix("/static/assets", http.FileServer(http.Dir("static/assets"))))
} else {
    log.Println("Serving embedded static files...")
    s.router.PathPrefix("/static/").
        Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Assets))))
}
```
why bother embedding? much harder to deploy incorrectly

For HTML
in debug parse during runtime so edits are picked up
for prod parse ahead of time to detect errors early

if not debug parse early
when template is executed already parsed?

else
template is nil and will be parsed everytime?

what if css and javascript was added to base.html instead of fetched as a separate file?