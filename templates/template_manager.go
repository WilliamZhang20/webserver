package templates

import (
	"bytes"
	"html/template"
	"log"
)

var tmpl *template.Template

func LoadTemplates() {
	var err error
	tmpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}
}

func ExecuteTemplate(name string, data interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
