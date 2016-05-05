package colorart

import "testing"

func TestColorString(t *testing.T) {
	c := Color{0.2, 0.4, 0.6, true}
	str := c.String()
	answer := "#336699"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}

	str = BlackColor.String()
	answer = "#000000"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}

	str = WhiteColor.String()
	answer = "#ffffff"
	if str != answer {
		t.Errorf("String conversion should be %s, not %s", answer, str)
	}
}

func TestColorRGBToColor(t *testing.T) {
	r := rgb{51, 102, 153}
	c := rgbToColor(r)

	if c.R != 0.2 || c.G != 0.4 || c.B != 0.6 {
		t.Errorf("Bad conversion from RGB(%d,%d,%d) to Color(%0.2f,%0.2f,%0.2f)", r[0], r[1], r[2], c.R, c.G, c.B)
	}
}

func TestColorIsBlackOrWhite(t *testing.T) {
	if !BlackColor.isBlackOrWhite() {
		t.Error("Black color is not BlackOrWhite")
	}

	if !WhiteColor.isBlackOrWhite() {
		t.Error("White color is not BlackOrWhite")
	}

	lightGray := Color{0.95, 0.95, 0.95, true}
	if !lightGray.isBlackOrWhite() {
		t.Error("LightGray color is not BlackOrWhite")
	}

	darkGray := Color{0.05, 0.05, 0.05, true}
	if !darkGray.isBlackOrWhite() {
		t.Error("DarkGray color is not BlackOrWhite")
	}

	red := Color{1, 0, 0, true}
	if red.isBlackOrWhite() {
		t.Error("Red color is BlackOrWhite")
	}
}

func TestColorIsDarkColor(t *testing.T) {
	lightGray := Color{0.95, 0.95, 0.95, true}
	if lightGray.isDarkColor() {
		t.Error("LightGray color is dark color")
	}

	darkGray := Color{0.05, 0.05, 0.05, true}
	if !darkGray.isDarkColor() {
		t.Error("DarkGray color is not dark color")
	}
}

func TestColorIsDistinctColor(t *testing.T) {
	red := Color{1, 0, 0, true}
	blue := Color{0, 0, 1, true}
	if !red.isDistinctColor(blue) {
		t.Error("Red is not distinct from blue")
	}

	almost := Color{1, 0.1, 0.1, true}
	if red.isDistinctColor(almost) {
		t.Error("Red is distinct from almost red")
	}
}

func TestColorIsContrastingColor(t *testing.T) {
	lightGray := Color{0.95, 0.95, 0.95, true}
	darkGray := Color{0.05, 0.05, 0.05, true}
	if lightGray.isContrastingColor(WhiteColor) {
		t.Error("LightGray is contrasting to white")
	}

	if !lightGray.isContrastingColor(darkGray) {
		t.Error("LightGray is not contrasting to dark gray")
	}
}

func TestColorMinSaturation(t *testing.T) {
	//c := Color{0.9, 0.8, 0.6, true}
	n := WhiteColor.colorWithMinimumSaturation(0.8)

	if n.G > 0.5 || n.B > 0.5 {
		t.Errorf("Color not desaturated enough from white (%0.2f,%0.2f,%0.2f)", n.R, n.G, n.B)
	}
}

func TestColorHSV(t *testing.T) {
	red := Color{1, 0, 0, true}
	h, s, v := red.hsv()
	if h != 0.0 || s != 1.0 || v != 1.0 {
		t.Errorf("Bad HSV conversion from Red (%f, %f, %f)", h, s, v)
	}

	blue := Color{0, 0, 1, true}
	h, s, v = blue.hsv()
	if h < 0.66 || h > 0.67 || s != 1.0 || v != 1.0 {
		t.Errorf("Bad HSV conversion from Blue (%f, %f, %f)", h, s, v)
	}

	green := Color{0, 1, 0, true}
	h, s, v = green.hsv()
	if h < 0.33 || h > 0.34 || s != 1.0 || v != 1.0 {
		t.Errorf("Bad HSV conversion from Green (%f, %f, %f)", h, s, v)
	}
}
