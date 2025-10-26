package util

import "errors"

// ErrorDuplicateKey is used when a unique constraint (e.g., email) is violated in the database.
var ErrorDuplicateKey = errors.New("duplicate Key Error")

// ErrorAuthenticated is used when authentication fails or the user is not logged in.
var ErrorAuthenticated = errors.New("Unauthenticated")
