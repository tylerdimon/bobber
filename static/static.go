package static

import (
	"bytes"
	"embed"
	"github.com/Masterminds/sprig/v3"
	"github.com/tylerdimon/bobber"
	"html/template"
	"io"
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
const endpointDetailPath = "html/endpoint-detail.html"
const actionButtonPath = "html/components/action-button.html"
const notFoundPath = "html/not-found.html"

type Executable interface {
	Execute(io.Writer, any) error
	ExecuteTemplate(io.Writer, string, any) error
}

var IndexTemplate Executable
var ConfigTemplate Executable
var RequestTemplate Executable
var RequestDetailTemplate Executable
var NamespaceAddTemplate Executable
var EndpointDetailTemplate Executable
var NotFoundTemplate Executable

type Reloader struct {
	templateGenerator func() *template.Template
}

func (r *Reloader) Execute(wr io.Writer, data any) error {
	htmlTemplate := r.templateGenerator()
	return htmlTemplate.Execute(wr, data)
}

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

	EndpointDetailTemplate, err = template.ParseFS(html, baseTemplatePath, endpointDetailPath)
	if err != nil {
		log.Fatal(err)
	}

	NotFoundTemplate, err = template.ParseFS(html, baseTemplatePath, notFoundPath)
	if err != nil {
		log.Fatal(err)
	}

}

// GetRequestHTML used to generate HTML pushed over websocket
func GetRequestHTML(request *bobber.Request) ([]byte, error) {
	var buf bytes.Buffer
	err := RequestTemplate.ExecuteTemplate(&buf, "request", request)
	if err != nil {
		log.Printf("error getting requst HTML to send over websocket: %v", err)
	}
	return buf.Bytes(), nil
}
