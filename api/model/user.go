package model

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        int        `json:"id"`
	Fullname  string     `json:"fullname" validate:"required"`
	Email     string     `json:"email" validate:"required,email"`
	Password  string     `json:"-" validate:"required"`
	Phone     string     `json:"phone" validate:"required,e164"`
	Active    bool       `json:"active"`
	IsAdmin   bool       `json:"is_admin"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type Register struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required,e164"`
}

type PasswordUpdate struct {
	ID         int    `json:"id" validate:"omitempty"` // This field is required if the administrator wants to update a user.
	Password   string `json:"password" validate:"required,min=6"`
	RePassword string `json:"re-password" validate:"required,min=6"`
}

type UserRequest struct {
	ID    int    `json:"id" param:"id"`
	Name  string `json:"name"`
	Role  string `json:"role" query:"role"`
	Token string `header:"Authorization"`
}

func (r *UserRequest) ParseBody(body io.Reader) error {
	return json.NewDecoder(body).Decode(r)
}

func (r *UserRequest) ParseParams(ctx *chi.Context) error {
	if id := ctx.URLParam("id"); id != "" {
		parsedID, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		r.ID = parsedID
	}
	return nil
}

func (r *UserRequest) ParseQuery(query url.Values) error {
	if role := query.Get("role"); role != "" {
		r.Role = role
	}
	return nil
}

func (r *UserRequest) ParseHeader(header http.Header) error {
	if token := header.Get("Authorization"); token != "" {
		r.Token = token
	}
	return nil
}
