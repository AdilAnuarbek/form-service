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
		Profile   Template
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
	user, err := h.getUser(r)
	if err != nil {
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
	type Attribute struct {
		attribute_name  string
		attribute_value string
	}
	var data struct {
		FormSTR    string
		Attributes []Attribute
	}
	data.FormSTR = chi.URLParam(r, "formSTR")

	// TODO: Check if user owns the form

	form_data, err := h.FormService.GetFormData(data.FormSTR)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, attributes := range form_data {
		data.Attributes = append(data.Attributes, Attribute{
			attribute_name:  attributes.Name,
			attribute_value: attributes.Value,
		})
	}

	h.Templates.ViewForm.Execute(w, r, data)
}

func (h Handlers) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	type Forms struct {
		FormName string
		FormSTR  string
	}
	var data struct {
		Forms []Forms
	}

	user, err := h.getUser(r)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	formNames, formSTR, err := h.FormService.GetForms(user.ID)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for i := 0; i < len(formNames); i++ {
		data.Forms = append(data.Forms, Forms{
			FormName: formNames[i],
			FormSTR:  formSTR[i],
		})
	}

	h.Templates.Profile.Execute(w, r, data)
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

// getEmail returns a user through context
func (h Handlers) getUser(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(ContextKeyUser).(*models.User)
	if !ok {
		return nil, fmt.Errorf("error no user")
	}
	return user, nil
}
