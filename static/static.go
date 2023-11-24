package static

import (
	"embed"
	"html/template"
	"log"
)

//go:embed assets
var Assets embed.FS

//go:embed html/index.html
var index embed.FS

var IndexTemplate *template.Template

func ParseHTML() {
	var err error
	IndexTemplate, err = template.ParseFS(index, "html/index.html")
	if err != nil {
		log.Fatal(err)
	}
}
