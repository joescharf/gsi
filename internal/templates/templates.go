package templates

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed files/*
var templateFS embed.FS

// Data holds the variables available to all templates.
type Data struct {
	ProjectName  string
	GoModulePath string
}

// Render executes the named template with the given data and returns the result.
func Render(name string, data Data) (string, error) {
	raw, err := templateFS.ReadFile("files/" + name)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(name).Parse(string(raw))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
