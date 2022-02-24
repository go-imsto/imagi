package image

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"

	"github.com/liut/jpegquality"
)

// consts
const (
	MinJPEGQuality = jpeg.DefaultQuality // 75
	MinWebpQuality = 80
)

// consts
const (
	FormatGIF  = "gif"
	FormatJPEG = "jpeg"
	FormatPNG  = "png"
	FormatWEBP = "webp"
)

var (
	mtypes = map[string]string{
		FormatGIF:  "image/gif",
		FormatJPEG: "image/jpeg",
		FormatPNG:  "image/png",
		FormatWEBP: "image/webp",
	}
)

// Imager ...
type Imager interface {
	SaveTo(w io.Writer, opt WriteOption) error
	ThumbnailTo(w io.Writer, topt ThumbOption) error
}

// Image ...
type Image struct {
	*Attr
	Format string
	m      image.Image
	rs     io.ReadSeeker
	rn     int // read length
}

// Open ...
func Open(rs io.ReadSeeker) (*Image, error) {

	cw := new(CountWriter)
	m, format, err := image.Decode(io.TeeReader(rs, cw))
	if err != nil {
		return nil, err
	}

	pt := m.Bounds().Max
	attr := NewAttr(uint(pt.X), uint(pt.Y), format)
	if mt, ok := mtypes[format]; ok {
		attr.Mime = mt
	}
	attr.Size = Size(cw.Len())
	if format == FormatJPEG {
		jr, err := jpegquality.New(rs)
		if err != nil {
			return nil, err
		}
		attr.Quality = Quality(jr.Quality())
	}
	return &Image{
		m:      m,
		Attr:   attr,
		Format: format,
		rs:     rs,
		rn:     cw.Len(),
	}, nil
}

// WriteOption ...
type WriteOption struct {
	Format  string
	Quality Quality
}

func (o *WriteOption) patch() {
	o.Format = PatchFormat(o.Format)
}

// SaveTo ...
func (im *Image) SaveTo(w io.Writer, opt WriteOption) error {
	if opt.Format == "" {
		opt.Format = im.Format
	}
	var buf bytes.Buffer
	n, err := SaveTo(&buf, im.m, opt)
	if err != nil {
		return err
	}
	var nn int64
	if n > im.rn {
		log.Printf("saved %d, im size %d", n, im.rn)
		im.rs.Seek(0, 0)
		nn, err = io.Copy(w, im.rs)
	} else {
		nn, err = io.Copy(w, &buf)
	}
	log.Printf("copied %d bytes", nn)
	return err
}

// SaveTo ...
func SaveTo(w io.Writer, m image.Image, opt WriteOption) (n int, err error) {
	cw := new(CountWriter)
	defer func() { n = cw.Len() }()
	w = io.MultiWriter(w, cw)
	opt.patch()
	switch opt.Format {
	case FormatJPEG:
		qlt := int(opt.Quality)
		if qlt == 0 {
			qlt = MinJPEGQuality
		}
		err = jpeg.Encode(w, m, &jpeg.Options{Quality: qlt})
		return
	case FormatGIF:
		err = gif.Encode(w, m, &gif.Options{
			NumColors: 256,
			Quantizer: nil,
			Drawer:    nil,
		})
		return
	case FormatPNG:
		err = png.Encode(w, m)
		return
	case FormatWEBP:
		qlt := int(opt.Quality)
		if qlt == 0 {
			qlt = MinWebpQuality
		}
		err = webpEncode(w, m, float32(qlt))
		return
	default:
		log.Printf("opt %v", opt)
		err = ErrUnsupportFormat
		return
	}

	return
}

// ThumbnailTo ...
func (im *Image) ThumbnailTo(w io.Writer, topt ThumbOption) error {
	if im.m == nil {
		return ErrEmptyImage
	}
	return ThumbnailImageTo(im.m, w, topt)
}
