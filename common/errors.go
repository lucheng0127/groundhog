package common

import (
	"github.com/pkg/errors"
)

// InitailizeError : Errors of initialize groundhog server
func InitailizeError(eMsg string) error {
	return errors.New(eMsg)
}

// IfaceSetupError : Errors of setup interface
func IfaceSetupError(eMsg string) error {
	return errors.New(eMsg)
}

// ServerError : Errors of UDP server
func ServerError(eMsg string) error {
	return errors.New(eMsg)
}
