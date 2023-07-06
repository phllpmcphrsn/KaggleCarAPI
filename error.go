package main

import "errors"

var errDbUsernameMissing = errors.New("database username not given or found")
var errDbPasswordMissing = errors.New("database password not given or found")
type APIError struct {
	ErrorCode    int
	ErrorMessage string
}