package views

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"bytes"
)

type Template struct {
	htmlTpl *template.Template
}

func ParseFileSys(pattern []string) (Template, error) {
	tpl, err := template.ParseFiles(pattern...)
	if err != nil {
		log.Printf("parsing template: %v", err)
		return Template{}, fmt.Errorf("could not parse the page %v", err)
	}

	return Template{
			htmlTpl: tpl,
		},
		nil

	//TODO: add FS later...
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return t
}

func (t Template) ExecuteTemplate(w http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Print("cloning: %w", err)
		http.Error(w, "Could not clone", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, "layout", data)
	if err != nil {
		log.Print("executing: %w", err)
		http.Error(w, "failed at executing", http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)

	//TODO: Later add "files" for csrf and other dynamic stuff
}
