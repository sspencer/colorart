package colorart

import "image"

const (
	// 1 to examine every pixel, 2 to skip every other pixel
	loopSkipper = 2

	// detune colors so colors within a few pixels of each other
	// map to the same color.  Makes algorthim much faster.
	// 0 is no change.
	// 1 divides by 2, multiplies by 2
	// 2 divides by 4, multiplies by 4
	// 3 divides by 8, multiplies by 8
	// Don't go much beyond 3...
	colorShifter = 2
)

// Colors is the return struct from Analyze.
type Colors struct {
	BackgroundColor, PrimaryColor, SecondaryColor, DetailColor Color
}

type colorArt struct {
	img *pixelGetter
}

// Analyze an image for its main colors.
func Analyze(img image.Image) Colors {
	c := &colorArt{}
	c.img = newPixelGetter(img)

	backgroundColor := c.findEdgeColor()
	primaryColor, secondaryColor, detailColor := c.findTextColors(backgroundColor)

	darkBackground := backgroundColor.isDarkColor()

	if !primaryColor.set {
		if darkBackground {
			primaryColor = WhiteColor
		} else {
			primaryColor = BlackColor
		}
	}

	if !secondaryColor.set {
		if darkBackground {
			secondaryColor = WhiteColor
		} else {
			secondaryColor = BlackColor
		}
	}

	if !detailColor.set {
		if darkBackground {
			detailColor = WhiteColor
		} else {
			detailColor = BlackColor
		}
	}

	return Colors{
		backgroundColor,
		primaryColor,
		secondaryColor,
		detailColor,
	}
}

func (c *colorArt) findTextColors(backgroundColor Color) (primaryColor, secondaryColor, detailColor Color) {
	b := c.img.imgBounds
	imageColors := parallelize(b.Min.Y, b.Max.Y, func(ch chan CountedSet, pmin, pmax int) {
		b := c.img.imgBounds
		colors := NewCountedSet(10000)
		for y := pmin; y < pmax; y += loopSkipper {
			for x := b.Min.X; x < b.Max.X; x += loopSkipper {
				colors.AddPixel(c.img.getPixel(x, y))
			}
		}

		ch <- colors
	})

	useDarkTextColor := !backgroundColor.isDarkColor()
	selectColors := NewCountedSet(5000)

	for key, cnt := range imageColors {
		// don't bother unless there's more than a few of the same color

		curColor := rgbToColor(key).colorWithMinimumSaturation(0.15)
		if curColor.isDarkColor() == useDarkTextColor {
			selectColors.AddCount(key, cnt)
		}

	}

	sortedColors := selectColors.SortedSet()
	for _, e := range sortedColors {
		curColor := rgbToColor(e.Color)
		if !primaryColor.set {
			if curColor.isContrastingColor(backgroundColor) {
				primaryColor = curColor
			}
		} else if !secondaryColor.set {
			if !primaryColor.isDistinctColor(curColor) || !curColor.isContrastingColor(backgroundColor) {
				continue
			}
			secondaryColor = curColor

		} else if !detailColor.set {
			if !secondaryColor.isDistinctColor(curColor) ||
				!primaryColor.isDistinctColor(curColor) ||
				!curColor.isContrastingColor(backgroundColor) {
				continue
			}
			detailColor = curColor
		}

		if primaryColor.set && secondaryColor.set && detailColor.set {
			break
		}
	}

	return
}

func (c *colorArt) findEdgeColor() Color {

	edgeColors := NewCountedSet(500)
	b := c.img.imgBounds
	x0 := b.Min.X
	x1 := b.Max.X - 1
	for y := b.Min.Y; y < b.Max.Y; y++ {
		edgeColors.AddPixel(c.img.getPixel(x0, y))
		edgeColors.AddPixel(c.img.getPixel(x1, y))
	}

	sortedColors := edgeColors.SortedSet()

	proposedEntry := sortedColors[0]
	proposedColor := rgbToColor(proposedEntry.Color)

	// try another color if edge is close to black or white
	if proposedColor.isBlackOrWhite() {
		for i, e := range sortedColors {
			if i == 0 {
				// first entry is already set as "proposedEntry"
				continue
			}

			nextProposedEntry := e
			// make sure second choice is 30% as common as first choice
			if float32(nextProposedEntry.Count)/float32(proposedEntry.Count) > 0.3 {
				nextProposedColor := rgbToColor(nextProposedEntry.Color)
				if !nextProposedColor.isBlackOrWhite() {
					proposedColor = nextProposedColor
					break
				}
			}
		}
	}

	return proposedColor
}
