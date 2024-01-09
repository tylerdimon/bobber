package static

import (
	"bytes"
	"embed"
	"github.com/Masterminds/sprig/v3"
	"github.com/tylerdimon/bobber"
	"html/template"
	"log"
)

//go:embed assets
var Assets embed.FS

//go:embed html*
var html embed.FS

const baseTemplatePath = "html/base.html"
const requestsIndexPath = "html/request-index.html"
const singleRequestPath = "html/request-short.html"
const requestDetailPath = "html/request-detail.html"
const configPath = "html/namespace-index.html"
const namespaceAddPath = "html/namespace-detail.html"
const endpointAddPath = "html/endpoint-add.html"
const actionButtonPath = "html/components/action-button.html"

var IndexTemplate *template.Template
var ConfigTemplate *template.Template
var RequestTemplate *template.Template
var RequestDetailTemplate *template.Template
var NamespaceAddTemplate *template.Template
var EndpointAddTemplate *template.Template

func ParseHTML() {
	var err error

	IndexTemplate, err = template.New("base.html").Funcs(sprig.FuncMap()).ParseFS(html, baseTemplatePath, requestsIndexPath, singleRequestPath, actionButtonPath)
	if err != nil {
		log.Fatal(err)
	}

	RequestTemplate, err = template.New("requests").Funcs(sprig.FuncMap()).ParseFS(html, singleRequestPath, actionButtonPath)
	if err != nil {
		log.Fatal(err)
	}

	RequestDetailTemplate, err = template.New("base.html").Funcs(sprig.FuncMap()).ParseFS(html, baseTemplatePath, requestDetailPath)
	if err != nil {
		log.Fatal(err)
	}

	ConfigTemplate, err = template.ParseFS(html, baseTemplatePath, configPath)
	if err != nil {
		log.Fatal(err)
	}

	NamespaceAddTemplate, err = template.ParseFS(html, baseTemplatePath, namespaceAddPath)
	if err != nil {
		log.Fatal(err)
	}

	EndpointAddTemplate, err = template.ParseFS(html, baseTemplatePath, endpointAddPath)
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
