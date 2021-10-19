//go:build !cgo
// +build !cgo

package image

import (
	"image"
	"io"

	_ "golang.org/x/image/webp" // ok
)

func webpEncode(w io.Writer, m image.Image, qlt float32) error {
	return ErrUnsupportFormat
}
