package common

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrInvalidInput = fmt.Errorf("invalid input")
)
