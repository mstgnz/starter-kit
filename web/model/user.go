package model

import (
	"time"
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
