package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/adilanuarbek/form-service/controllers"
	"github.com/adilanuarbek/form-service/models"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	tpl = tpl.Funcs(template.FuncMap{
		"userSignedin": func() *models.User {
			user, ok := r.Context().Value(controllers.ContextKeyUser).(*models.User)
			if !ok {
				return nil
			}
			return user
		},
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	if data == nil {
		var newdata struct {
			Email string
		}
		user, ok := r.Context().Value(controllers.ContextKeyUser).(*models.User)
		if ok {
			newdata.Email = user.Email
		}
		err = tpl.Execute(&buf, newdata)
	} else {
		err = tpl.Execute(&buf, data)
	}

	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(filepath.Base(patterns[0]))
	tpl = tpl.Funcs(template.FuncMap{
		"userSignedin": func() (template.HTML, error) {
			return "", fmt.Errorf("currentUser not implemented")
		},
	})
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}
	return Template{htmlTpl: tpl}, nil
}
