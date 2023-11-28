package static

import (
	"bytes"
	"embed"
	"github.com/tylerdimon/bobber"
	"html/template"
	"log"
)

//go:embed assets
var Assets embed.FS

//go:embed html*
var html embed.FS

const baseTemplatePath = "html/base.html"
const requestsIndexPath = "html/index.html"
const singleRequestPath = "html/request.html"
const configPath = "html/config.html"

var IndexTemplate *template.Template
var ConfigTemplate *template.Template
var RequestTemplate *template.Template

func ParseHTML() {
	var err error

	IndexTemplate, err = template.ParseFS(html, baseTemplatePath, requestsIndexPath, singleRequestPath)
	if err != nil {
		log.Fatal(err)
	}

	RequestTemplate, err = template.ParseFS(html, singleRequestPath)
	if err != nil {
		log.Fatal(err)
	}

	ConfigTemplate, err = template.ParseFS(html, baseTemplatePath, configPath)
	if err != nil {
		log.Fatal(err)
	}
}

func GetRequestHTML(request *bobber.Request) ([]byte, error) {
	var buf bytes.Buffer
	err := RequestTemplate.ExecuteTemplate(&buf, "request", request)
	if err != nil {
		log.Printf("error getting requst HTML to send over websocket: %v", err)
	}
	return buf.Bytes(), nil
}
