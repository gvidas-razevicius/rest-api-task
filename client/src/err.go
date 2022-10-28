package client

import (
	"errors"
	"fmt"

	color "github.com/fatih/color"
)

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

var (
	ErrNotFound  error = errors.New(red("one or more objects were not found!"))
	ErrForbidden error = errors.New(red("object already exists!"))
)

type ErrBadStatus struct {
	StatusMsg string
}

type ErrTypeInvalid struct {
	ArgName string
	TypeErr error
}

type ErrMakeRequest struct {
	RequestErr error
}

type ErrEncodeJson struct {
	JsonError error
}

type ErrDecodeJson struct {
	JsonError error
}

type ErrServerNegative struct {
	NegativeMsg error
}

func (err ErrBadStatus) Error() string {
	return fmt.Sprintf("server returned: \n%s", err.StatusMsg)
}

func (err ErrMakeRequest) Error() string {
	return fmt.Sprintf("could not make request: %v", err.RequestErr)
}

func (err ErrEncodeJson) Error() string {
	return fmt.Sprintf("could not encode json data: %v", err.JsonError)
}

func (err ErrDecodeJson) Error() string {
	return fmt.Sprintf("could not decode results into json: %v", err.JsonError)
}

func (err ErrTypeInvalid) Error() string {
	return fmt.Sprintf(red("%s type is invalid: %v"), err.ArgName, err.TypeErr)
}

func (err ErrServerNegative) Error() string {
	return fmt.Sprintf("server could not process request: \n%s", err.NegativeMsg)
}
