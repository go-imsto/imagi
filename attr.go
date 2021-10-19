package image

// Dimension ...
type Dimension uint32

// Quality ...
type Quality uint8

// Size ...
type Size uint32

// Attr ...
type Attr struct {
	Width   Dimension `json:"width"`
	Height  Dimension `json:"height"`
	Quality Quality   `json:"qlt,omitempty"`  // Original quality
	Size    Size      `json:"size,omitempty"` // Original size
	Ext     string    `json:"ext"`            // file extension include dot
	Mime    string    `json:"mime,omitempty"` // content type
}

// ToMap ...
func (a Attr) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"width":  a.Width,
		"height": a.Height,
		"ext":    a.Ext,
		"mime":   a.Mime,
	}

	if a.Quality > 0 {
		m["qlt"] = a.Quality
	}
	return m
}

// FromMap ...
func (a *Attr) FromMap(m map[string]interface{}) {
	if m == nil {
		return
	}
	if a == nil {
		*a = Attr{}
	}
	if v, ok := m["width"]; ok {
		if vv, ok := v.(uint32); ok {
			a.Width = Dimension(vv)
		}
	}
	if v, ok := m["height"]; ok {
		if vv, ok := v.(uint32); ok {
			a.Height = Dimension(vv)
		}
	}
	if v, ok := m["ext"]; ok {
		if vv, ok := v.(string); ok {
			a.Ext = vv
		}
	}
	if v, ok := m["mime"]; ok {
		if vv, ok := v.(string); ok {
			a.Mime = vv
		}
	}
}

// NewAttr ...
func NewAttr(w, h uint, f string) *Attr {
	a := &Attr{
		Width:  Dimension(w),
		Height: Dimension(h),
		Ext:    Format2Ext(f),
	}
	return a
}

// Format2Ext ...
func Format2Ext(f string) string {
	if f == "jpeg" {
		return ".jpg"
	}
	return "." + f
}

// Ext2Format ...
func Ext2Format(s string) string {
	if s == "" {
		return s
	}
	if s[0] == '.' {
		s = s[1:]
	}
	if s == "jpg" {
		return "jpeg"
	}
	return s
}
