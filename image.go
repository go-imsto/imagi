package image

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"

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
	im, err := NewFromImage(m, cw.Len(), format)
	if err != nil {
		return nil, err
	}
	im.rs = rs
	if format == FormatJPEG {
		jr, err := jpegquality.New(rs)
		if err != nil {
			return nil, err
		}
		im.Quality = uint8(jr.Quality())
	}
	return im, nil
}

func NewFromImage(m image.Image, size int, format string) (*Image, error) {
	pt := m.Bounds().Max
	attr := NewAttr(uint(pt.X), uint(pt.Y), format)
	if mt, ok := mtypes[format]; ok {
		attr.Mime = mt
	}
	attr.Size = uint32(size)
	return &Image{
		m: m, Attr: attr, Format: format, rn: size,
	}, nil
}

// WriteOption ...
type WriteOption struct {
	Format  string
	Quality uint8

	ExtraWriter io.Writer // 额外的输出 一般用于hash计算
}

func (o *WriteOption) patch() {
	o.Format = PatchFormat(o.Format)
}

// SaveTo ...
func (im *Image) SaveTo(w io.Writer, opt *WriteOption) (int, error) {
	if opt == nil {
		opt = new(WriteOption)
	}
	if opt.Format == "" {
		opt.Format = im.Format
	}
	if !WebpEncodable && im.Format == FormatWEBP && im.rs != nil {
		_, _ = im.rs.Seek(0, 0)
		n, err := io.Copy(w, im.rs)
		return int(n), err
	}
	var buf bytes.Buffer
	err := SaveTo(&buf, im.m, opt)
	if err != nil {
		return 0, err
	}
	var nn int64
	if im.Format == opt.Format && buf.Len() > im.rn && im.rs != nil {
		slog.Debug("saved", "n", buf.Len(), "read length", im.rn)
		_, _ = im.rs.Seek(0, 0)
		nn, err = io.Copy(w, im.rs)
	} else {
		nn, err = io.Copy(w, &buf)
	}
	if err != nil {
		slog.Info("copy fail", "err", err, "bytes", nn)
	} else {
		slog.Debug("copied", "bytes", nn)
	}
	return int(nn), err
}

// SaveTo ...
func SaveTo(w io.Writer, m image.Image, opt *WriteOption) (err error) {
	if opt == nil {
		opt = new(WriteOption)
	}

	if opt.ExtraWriter != nil {
		w = io.MultiWriter(w, opt.ExtraWriter)
	}

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
		slog.Info("invalid format", "opt", opt)
		err = ErrUnsupportFormat
		return
	}
}

// ThumbnailTo ...
func (im *Image) ThumbnailTo(w io.Writer, topt *ThumbOption) error {
	if im.m == nil {
		return ErrEmptyImage
	}
	if topt.Format == "" {
		topt.Format = im.Format
	}
	return ThumbnailImageTo(im.m, w, topt)
}
