package image

// Attr ...
type Attr struct {
	Width   uint32 `json:"width"`
	Height  uint32 `json:"height"`
	Size    uint32 `json:"size,omitempty"` // Original size
	Quality uint8  `json:"qlt,omitempty"`  // Original quality
	Ext     string `json:"ext"`            // file extension include dot
	Mime    string `json:"mime,omitempty"` // content type
}

// ToMap ...
func (a Attr) ToMap() map[string]any {
	m := map[string]any{
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
func (a *Attr) FromMap(m map[string]any) {
	if m == nil {
		return
	}
	if a == nil {
		*a = Attr{}
	}
	if v, ok := m["width"]; ok {
		if vv, ok := v.(uint32); ok {
			a.Width = vv
		}
	}
	if v, ok := m["height"]; ok {
		if vv, ok := v.(uint32); ok {
			a.Height = vv
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
		Width:  uint32(w),
		Height: uint32(h),
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

// PatchFormat ...
func PatchFormat(s string) string {
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
