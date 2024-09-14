package models

import (
	"errors"
)

// Allow our handlers to have a common erro interface to manage
// that remains agnostic to the underlying datastore
var ErrNoRecord = errors.New("models: no matching record found")