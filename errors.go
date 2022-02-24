package image

import (
	"errors"
)

// vars of errors
var (
	ErrInvalidFormat   = errors.New("invalid image format")
	ErrUnsupportFormat = errors.New("unsupported image format")
	ErrOrigTooSmall    = errors.New("original image too small")
	ErrEmptyImage      = errors.New("image is empty")
)
