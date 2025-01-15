package controllers

import (
	"fmt"
	"net/http"

	"github.com/adilanuarbek/form-service/models"
	"github.com/go-chi/chi/v5"
)

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}

type Handlers struct {
	Templates struct {
		Index     Template
		Contact   Template
		Dashboard Template
		ViewForm  Template
	}
	FormService *models.FormService
}

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func (h Handlers) PostDashboardHandler(w http.ResponseWriter, r *http.Request) {
	formName := r.FormValue("form-name")
	formSTR := models.RandomString(8)
	for h.FormService.CheckFormStr(formSTR) { // ensuring that formSTR does not repeat
		formSTR = models.RandomString(8)
	}
	user, ok := r.Context().Value(ContextKeyUser).(*models.User)
	if !ok {
		http.Error(w, "User not logged in", http.StatusInternalServerError)
		return
	}

	nf := models.NewForm{
		FormName: formName,
		FormSTR:  formSTR,
		UserID:   user.ID,
	}

	form, err := h.FormService.CreateForm(nf)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// http.Redirect(w, r, "/dashboard", http.StatusFound)
	http.Redirect(w, r, fmt.Sprintf("/%v", form.FormSTR), http.StatusFound)
}

func (h Handlers) ViewFormData(w http.ResponseWriter, r *http.Request) {
	var data struct {
		formSTR string
		data    map[string]interface{}
	}
	data.formSTR = chi.URLParam(r, "formSTR")

	h.Templates.ViewForm.Execute(w, r, data)
}

func (h Handlers) ACTUALSERVICE(w http.ResponseWriter, r *http.Request) { // TODO: Change name
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Iterate through all form keys and values
	for key, values := range r.Form {
		// `values` is a slice of strings
		fmt.Printf("Key: %s, Values: %v\n", key, values)
	}

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}
