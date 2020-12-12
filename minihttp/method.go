package minihttp

import (
	"errors"
	"strings"
)

var (
	//ErrNotImplementedMethod for unkown methods
	ErrNotImplementedMethod = errors.New("Not implemented method")
)

var supportedMethods = [...]string{"GET", "POST"}

func isValidMethod(method string) bool {
	for _, m := range supportedMethods {
		if strings.ToUpper(method) == m {
			return true
		}
	}
	return false
}
