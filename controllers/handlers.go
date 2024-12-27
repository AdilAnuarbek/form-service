package controllers

import (
	"net/http"
)

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}

type Handlers struct {
	Templates struct {
		Index   Template
		Contact Template
	}
}

func (h Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {

	h.Templates.Index.Execute(w, r, nil)
}

func (h Handlers) ContactHandler(w http.ResponseWriter, r *http.Request) {
	h.Templates.Contact.Execute(w, r, nil)
}
