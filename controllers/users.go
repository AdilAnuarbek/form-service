package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/adilanuarbek/form-service/models"
	"github.com/gorilla/sessions"
)

type contextKey string

const ContextKeyUser = contextKey("user")

type UserMiddleware struct {
	UserService *models.UserService
	Session     *sessions.CookieStore
}

type Users struct {
	Templates struct {
		Signup Template
	}
	UserService *models.UserService
	// SessionService *models.SessionService
	Session *sessions.CookieStore
}

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) { // GET

	u.Templates.Signup.Execute(w, r, nil)
}

func (u Users) PostSignUp(w http.ResponseWriter, r *http.Request) { // POST
	email := r.FormValue("email")
	password := r.FormValue("password")
	nu := models.NewUser{
		Email:    email,
		Password: password,
	}

	user, err := u.UserService.CreateUser(nu)
	if err != nil {
		fmt.Printf("PostSignUp: %v", err)
		http.Error(w, "Something went wrong. Check the console", http.StatusInternalServerError)
		return
	}

	session, _ := u.Session.Get(r, "form-service")
	session.Values["userID"] = user.ID

	err = sessions.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound) // TODO: change to "/"
}

func (umw UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := umw.Session.Get(r, "form-service")
		userID, ok := session.Values["userID"].(int)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.UserService.FindUser(userID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ContextKeyUser, user)
		fmt.Println("HELLO")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(ContextKeyUser).(*models.User)
		if !ok {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
