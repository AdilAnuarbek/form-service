package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
)

type FormData struct {
	Name  string
	Value string
}

type NewForm struct {
	FormName string
	FormSTR  string
	UserID   int
}

type Form struct {
	ID       int
	FormName string
	FormSTR  string
	UserID   int
}

type FormService struct {
	DB *sql.DB
}

func (fs *FormService) CreateForm(nf NewForm) (*Form, error) {
	form := Form{
		FormName: nf.FormName,
		FormSTR:  nf.FormSTR,
		UserID:   nf.UserID,
	}

	row := fs.DB.QueryRow(`INSERT INTO forms (user_id, form_str, form_name)
	VALUES ($1, $2, $3) RETURNING id`, nf.UserID, nf.FormSTR, nf.FormName)
	err := row.Scan(&form.ID)
	if err != nil {
		return nil, fmt.Errorf("models: failed to insert form: %w", err)
	}

	return &form, nil
}

// GetFormData returns data corresponding to a specific form with formSTR
func (fs *FormService) GetFormData(formSTR string) ([]FormData, error) {
	return nil, nil
}

// GetForms returns all forms created by a user with userID
func (fs *FormService) GetForms(userID int) ([]string, []string, error) {
	rows, err := fs.DB.Query(`SELECT form_name, form_str FROM forms WHERE user_id = $1`, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("models: failed to query forms: %w", err)
	}
	var formNames, formSTR []string

	for rows.Next() {
		var form_name, form_str string
		if err = rows.Scan(&form_name, &form_str); err != nil {
			return nil, nil, fmt.Errorf("models: failed to scan: %w", err)
		}
		formNames = append(formNames, form_name)
		formSTR = append(formSTR, form_str)
	}
	return formNames, formSTR, nil
}

// CheckFormStr returns false if the formSTR is unique
func (fs *FormService) CheckFormStr(formSTR string) bool {
	var exists bool
	err := fs.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM forms WHERE form_str = $1)`,
		formSTR).Scan(&exists)
	return err == sql.ErrNoRows
}

func (fs *FormService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
