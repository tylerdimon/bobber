package static

import (
	"embed"
	"github.com/tylerdimon/bobber"
	"html/template"
	"log"
	"os"
)

//go:embed assets
var Assets embed.FS

//go:embed html*
var html embed.FS

var IndexTemplate *template.Template

func ParseHTML() {
	var err error

	IndexTemplate, err = template.ParseFS(html, "html/base.html", "html/index.html")
	if err != nil {
		log.Fatal(err)
	}
}

type PageData struct {
	Title string
	Data  any
}

func ExecIndexTemplate(data []bobber.Request) error {
	pageData := PageData{
		Title: "Requests",
		Data:  data,
	}
	err := IndexTemplate.Execute(os.Stdout, pageData)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
