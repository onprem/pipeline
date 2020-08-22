package main

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func getIndexHTML() ([]byte, error) {
	src, err := ioutil.ReadFile("README.md")
	if err != nil {
		return nil, err
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(&buf)
	if err != nil {
		return nil, err
	}

	tpl := template.Must(template.New("index.html").ParseFiles("templates/index.html"))
	data := struct{ Body string }{Body: string(b)}

	var idx bytes.Buffer
	err = tpl.Execute(&idx, data)
	if err != nil {
		return nil, err
	}

	index, err := ioutil.ReadAll(&idx)
	if err != nil {
		return nil, err
	}

	return index, nil
}
