package api

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var Templates = map[string]*template.Template{}

func LoadTemplates(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t, err := template.ParseFiles(path)

			if err != nil {
				return err
			}

			rel, err := filepath.Rel(dir, path)

			if err != nil {
				return err
			}

			Templates[rel] = t
		}

		return nil
	})
}

func RenderTemplate(w http.ResponseWriter, path string, params interface{}) *Error {
	t, ok := Templates[path]

	if !ok {
		return ServerErrorf("no such template: %s", path)
	}

	if err := t.Execute(w, params); err != nil {
		return ServerError(err)
	}

	return nil
}
