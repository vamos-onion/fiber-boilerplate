package defs

import (
	"errors"
	"net/http"
	"strings"
	"syscall"
)

const (
	KeyStore string = "fiber-boilerplate/#/keyStore"
)

// errors
var (
	ErrBadRequest      = NewError(http.StatusBadRequest)
	ErrFault           = NewError(syscall.EFAULT)
	ErrInvalid         = NewError(syscall.EINVAL)
	ErrNoContent       = NewError(http.StatusNoContent)
	ErrNotDirectory    = NewError(syscall.ENOTDIR)
	ErrNotFound        = NewError(http.StatusNotFound)
	ErrNotImplemented  = NewError(http.StatusNotImplemented)
	ErrUnauthorized    = NewError(http.StatusUnauthorized)
	ErrConflict        = NewError(http.StatusConflict)
	ErrForbidden       = NewError(http.StatusForbidden)
	ErrInternalServer  = NewError(http.StatusInternalServerError)
	ErrUpgradeRequired = NewError(http.StatusUpgradeRequired)
	ErrTimeout         = NewError(http.StatusRequestTimeout)
)

// NewError :
func NewError(arg interface{}) error {
	var str string

	switch arg.(type) {
	case syscall.Errno:
		str = arg.(syscall.Errno).Error()
	case string:
		str = arg.(string)
	case int:
		str = http.StatusText(arg.(int))
	default:
		panic(syscall.EINVAL.Error())
	}

	return errors.New(strings.ToUpper(str))
}
