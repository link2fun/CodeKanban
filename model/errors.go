package model

import "errors"

// ErrDBNotInitialized indicates the SQL layer has not been prepared.
var ErrDBNotInitialized = errors.New("database is not initialized")
