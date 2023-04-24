//go:build !cgo
// +build !cgo

package image

import (
	"errors"
	"image"
	"io"

	_ "golang.org/x/image/webp" // ok
)

const (
	WebpEncodable = false
)

var (
	ErrUnsupportEncodeWebP = errors.New("unsupported encode webP")
)

func webpEncode(w io.Writer, m image.Image, qlt float32) error {
	return ErrUnsupportEncodeWebP
}
