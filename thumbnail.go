package image

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/nfnt/resize"
)

// ThumbOption 缩图选项
type ThumbOption struct {
	Width, Height       uint // 宽和高
	MaxWidth, MaxHeight uint // 最大宽和高
	IsFit               bool // 是否保持比例
	IsCrop              bool // 是否裁切
	CropX, CropY        int  // 裁切位置

	ctWidth, ctHeight uint // for crop temporary

	WriteOption
}

func (topt ThumbOption) String() string {
	return fmt.Sprintf("%dx%d q%d %v %v", topt.Width, topt.Height, topt.Quality, topt.IsFit, topt.IsCrop)
}

func (topt *ThumbOption) calc(ow, oh uint) error {
	if topt.Width >= ow && topt.Height >= oh {
		return fmt.Errorf("%dx%d is too big, orig is %dx%d", topt.Width, topt.Height, ow, oh)
	}

	if topt.IsFit {
		if topt.IsCrop {
			ratioX := float32(topt.Width) / float32(ow)
			ratioY := float32(topt.Height) / float32(oh)

			if ratioX > ratioY {
				topt.ctWidth = topt.Width
				topt.ctHeight = uint(ratioX * float32(oh))
			} else {
				topt.ctHeight = topt.Height
				topt.ctWidth = uint(ratioY * float32(ow))
			}
			// :resize

			if topt.ctWidth == topt.Width && topt.ctHeight == topt.Height {
				return nil
			}

			topt.CropX = int(float32(topt.ctWidth-topt.Width) / 2)
			topt.CropY = int(float32(topt.ctHeight-topt.Height) / 2)

			// slog.Debug("opt crop", "cropX", topt.CropX, "cropY", topt.CropY)

		} else {

			rel := float32(ow) / float32(oh)
			if topt.MaxWidth > 0 && topt.MaxWidth <= ow {
				topt.Width = topt.MaxWidth
				topt.Height = uint(float32(topt.Width) / rel)
			} else if topt.MaxHeight > 0 && topt.MaxHeight <= oh {
				topt.Height = topt.MaxHeight
				topt.Width = uint(float32(topt.Height) * rel)
			} else {
				bounds := float32(topt.Width) / float32(topt.Height)
				if rel >= bounds {
					topt.Height = uint(float32(topt.Width) / rel)
				} else {
					topt.Width = uint(float32(topt.Height) * rel)
				}
			}
		}
	}
	return nil
}

func (topt *ThumbOption) GetHeight() uint {
	if topt.IsFit && topt.IsCrop {
		return topt.ctHeight
	}
	return topt.Height
}

func (topt *ThumbOption) GetWidth() uint {
	if topt.IsFit && topt.IsCrop {
		return topt.ctWidth
	}
	return topt.Width
}

// ThumbnailImage ...
func ThumbnailImage(img image.Image, topt *ThumbOption) (image.Image, error) {

	ob := img.Bounds()
	ow := uint(ob.Dx())
	oh := uint(ob.Dy())

	if ow <= topt.Width && oh <= topt.Height {
		slog.Debug("ThumbnailImage", "ow", ow, "oh", oh, "w", topt.Width, "h", topt.Height)
		return img, nil
	}

	err := topt.calc(ow, oh)
	if err != nil {
		return nil, err
	}
	// slog.Debug("ThumbnailImage", "topt", topt)
	if topt.IsFit {
		if topt.IsCrop {
			buf := resize.Resize(topt.ctWidth, topt.ctHeight, img, resize.Bicubic)
			dst := image.NewRGBA(image.Rect(0, 0, int(topt.Width), int(topt.Height)))
			pt := image.Point{topt.CropX, topt.CropY}
			draw.Draw(dst, dst.Bounds(), buf, pt, draw.Src)
			return dst, nil
		}
	}
	m := resize.Resize(topt.Width, topt.Height, img, resize.Bicubic)
	return m, nil
}

// Thumbnail ...
func Thumbnail(r io.Reader, w io.Writer, topt *ThumbOption) error {
	var err error
	im, format, err := image.Decode(r)
	if err != nil {
		slog.Info("Thumbnail image decode fail", "err", err)
		return err
	}
	if topt.Format == "" {
		topt.Format = format
	}

	err = ThumbnailImageTo(im, w, topt)
	if err != nil {
		if err == ErrOrigTooSmall {
			if rr, ok := r.(io.Seeker); ok {
				_, _ = rr.Seek(0, 0)
			}
			var written int64
			written, err = io.Copy(w, r)
			if err == nil {
				slog.Debug("copied", "n", written)
				return nil
			}
			slog.Info("copy fail", "err", err)
		}
	}
	return err
}

// ThumbnailImageTo ...
func ThumbnailImageTo(im image.Image, w io.Writer, topt *ThumbOption) error {
	m, err := ThumbnailImage(im, topt)
	if err != nil {
		return err
	}

	opt := &topt.WriteOption
	err = SaveTo(w, m, opt)
	if err != nil {
		slog.Info("save to", "err", err)
		return err
	}

	return nil
}

// ThumbnailFile ...
func ThumbnailFile(src, dest string, topt *ThumbOption) (err error) {
	var in *os.File
	in, err = os.Open(src)
	if err != nil {
		slog.Info("open fail", "err", err)
		return
	}
	defer in.Close()

	dir := path.Dir(dest)
	err = os.MkdirAll(dir, os.FileMode(0755))
	if err != nil {
		return
	}

	var out *os.File
	out, err = os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		slog.Info("openfile fail", "err", err)
		return
	}
	defer out.Close()

	err = Thumbnail(in, out, topt)

	return
}
