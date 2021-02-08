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

// IPPoolEmptyError : Errors of ip pool no enough ip
func IPPoolEmptyError(eMsg string) error {
	return errors.New(eMsg)
}

// LaunchServerError : Errors of launch udp server
func LaunchServerError(eMsg string) error {
	return errors.New(eMsg)
}

// UDPConnectionError : Errors of UDP connection
func UDPConnectionError(eMsg string) error {
	return errors.New(eMsg)
}

// IfaceReadError : Errors of read data from interface
func IfaceReadError(eMsg string) error {
	return errors.New(eMsg)
}
