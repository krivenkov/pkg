package es

import "fmt"

var (
	ErrNotFound  = fmt.Errorf("obj not found")
	ErrMultiRows = fmt.Errorf("multiple rows in result set")
)
