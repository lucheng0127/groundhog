package common

import (
	"github.com/pkg/errors"
)

// InitailizeError return InitializeError
func InitailizeError(eMsg string) error {
	return errors.New(eMsg)
}
