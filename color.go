package colorart

import (
	"fmt"
	"math"
)

var (
	// BlackColor convenience
	BlackColor = Color{0, 0, 0, true}

	// WhiteColor convenience
	WhiteColor = Color{1, 1, 1, true}
)

// Color holds rgb components of a color
type Color struct {
	R, G, B float64
	set     bool
}

// String returns a HTML hex code string for the color (#9a45bc)
func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", uint8(c.R*255), uint8(c.G*255), uint8(c.B*255))
}

// RGBAToColor converts 16bit RGBA (0-65535) color into a Color (0.0 < rgb component < 1.0)
func RGBAToColor(r, g, b, a uint32) Color {
	fa := float64(a)
	return Color{float64(r) / fa, float64(g) / fa, float64(b) / fa, true}
}

// StringToColor converts [3]byte array [0xff, 0x00, 0x33]) to Color
func rgbToColor(c rgb) Color {
	return Color{float64(c[0]) / 255.0, float64(c[1]) / 255.0, float64(c[2]) / 255.0, true}
}

// IsBlackOrWhite returns true if the color is within about 90% or black or white
func (c Color) IsBlackOrWhite() bool {
	return (c.R > 0.91 && c.G > 0.91 && c.B >= 0.91) || (c.R < 0.09 && c.G < 0.09 && c.B < 0.09)
}

// IsDarkColor returns true for dark colors
func (c Color) IsDarkColor() bool {

	lum := 0.2126*c.R + 0.7152*c.G + 0.0722*c.B

	return lum < 0.5
}

// IsDistinctColor uses a minimum threshold to determine if two colors are distinct
func (c Color) IsDistinctColor(d Color) bool {
	threshold := 0.25
	if math.Abs(c.R-d.R) > threshold || math.Abs(c.G-d.G) > threshold || math.Abs(c.B-d.B) > threshold {

		threshold = 0.03
		// check for grays, prevent multiple gray colors
		if math.Abs(c.R-c.G) < threshold && math.Abs(c.R-c.B) < threshold {
			if math.Abs(d.R-d.G) < threshold && math.Abs(d.R-d.B) < threshold {
				return false
			}
		}

		return true
	}

	return false
}

// IsContrastingColor determines if two colors are contrasting
func (c Color) IsContrastingColor(d Color) bool {

	bLum := 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
	fLum := 0.2126*d.R + 0.7152*d.G + 0.0722*d.B

	contrast := 0.0

	if bLum > fLum {
		contrast = (bLum + 0.05) / (fLum + 0.05)
	} else {
		contrast = (fLum + 0.05) / (bLum + 0.05)
	}

	return contrast > 1.6
}

// ColorWithMinimumSaturation tries to return a less saturated color
func (c Color) ColorWithMinimumSaturation(minSaturation float64) Color {

	h, s, v := c.HSV()
	if s < minSaturation {
		return HSVToColor(h, minSaturation, v)
	}

	return c
}

// HSV converts color into HSV
// Ported from http://goo.gl/Vg1h9
func (c Color) HSV() (h, s, v float64) {

	max := math.Max(math.Max(c.R, c.G), c.B)
	min := math.Min(math.Min(c.R, c.G), c.B)
	d := max - min
	s, v = 0, max
	if max > 0 {
		s = d / max
	}
	if max == min {
		// Achromatic.
		h = 0
	} else {
		// Chromatic.
		switch max {
		case c.R:
			h = (c.G - c.B) / d
			if c.G < c.B {
				h += 6
			}
		case c.G:
			h = (c.B-c.R)/d + 2
		case c.B:
			h = (c.R-c.G)/d + 4
		}
		h /= 6
	}
	return
}

// HSVToColor converts an HSV triple to a (RGB) color.
//
// Ported from http://goo.gl/Vg1h9
func HSVToColor(h, s, v float64) Color {
	var r, g, b float64
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1.0 - s)
	q := v * (1.0 - f*s)
	t := v * (1.0 - (1.0-f)*s)
	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return Color{r, g, b, true}
}
