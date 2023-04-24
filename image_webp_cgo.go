//go:build cgo
// +build cgo

package image

import (
	"image"
	"io"

	"github.com/chai2010/webp"
)

const (
	WebpEncodable = true
)

func webpEncode(w io.Writer, m image.Image, qlt float32) error {
	return webp.Encode(w, m, &webp.Options{Quality: qlt})
}
