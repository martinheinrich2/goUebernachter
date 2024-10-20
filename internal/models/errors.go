package models

import (
	"errors"
)

// List various error messages
var (
	// ErrNoRecord error is used if no matching record is found in the database.
	ErrNoRecord = errors.New("models: no matching record found")

	// The ErrInvalidCredentials error is used if a user tries to login
	// with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// The ErrDuplicateEmail error will be used if a user tries to signup
	// with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")

	FailedValidation = errors.New("models: failed validation")
)
