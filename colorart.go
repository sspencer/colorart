package colorart

import "image"

// Colors is the return struct from Analyze.
type Colors struct {
	BackgroundColor, PrimaryColor, SecondaryColor, DetailColor Color
}

type colorArt struct {
	img *pixelGetter
}

// Analyze an image for its main colors.
// 1 to examine every pixel, 2 to skip every other pixel

func Analyze(img image.Image, colorShift int, loopSkip int) Colors {
	c := &colorArt{}
	c.img = newPixelGetter(img)

	colorShift = minMax(colorShift, 0, 8)
	loopSkip = minMax(loopSkip, 1, 8)

	backgroundColor := c.findEdgeColor(colorShift)
	primaryColor, secondaryColor, detailColor := c.findImageColors(backgroundColor, colorShift, loopSkip)

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

func (c *colorArt) findImageColors(backgroundColor Color, colorShift int, loopSkip int) (primaryColor, secondaryColor, detailColor Color) {
	b := c.img.imgBounds
	imageColors := parallelize(b.Min.Y, b.Max.Y, func(ch chan countedSet, pmin, pmax int) {
		b := c.img.imgBounds
		colors := newCountedSet(10000, colorShift)
		for y := pmin; y < pmax; y += loopSkip {
			for x := b.Min.X; x < b.Max.X; x += loopSkip {
				colors.addPixel(c.img.getPixel(x, y))
			}
		}

		ch <- colors
	})

	useDarkTextColor := !backgroundColor.isDarkColor()
	selectColors := newCountedSet(5000, colorShift)

	for key, cnt := range imageColors.set {
		// don't bother unless there's more than a few of the same color

		curColor := rgbToColor(key).colorWithMinimumSaturation(0.15)
		if curColor.isDarkColor() == useDarkTextColor {
			selectColors.addCount(key, cnt)
		}

	}

	sortedColors := selectColors.sortedSet()
	for _, e := range sortedColors {
		curColor := rgbToColor(e.color)
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

func (c *colorArt) findEdgeColor(colorShift int) Color {

	edgeColors := newCountedSet(500, colorShift)
	b := c.img.imgBounds
	x0 := b.Min.X
	x1 := b.Max.X - 1
	for y := b.Min.Y; y < b.Max.Y; y++ {
		edgeColors.addPixel(c.img.getPixel(x0, y))
		edgeColors.addPixel(c.img.getPixel(x1, y))
	}

	sortedColors := edgeColors.sortedSet()

	proposedEntry := sortedColors[0]
	proposedColor := rgbToColor(proposedEntry.color)

	// try another color if edge is close to black or white
	if proposedColor.isBlackOrWhite() {
		for i, e := range sortedColors {
			if i == 0 {
				// first entry is already set as "proposedEntry"
				continue
			}

			nextProposedEntry := e
			// make sure second choice is 30% as common as first choice
			if float32(nextProposedEntry.count)/float32(proposedEntry.count) > 0.3 {
				nextProposedColor := rgbToColor(nextProposedEntry.color)
				if !nextProposedColor.isBlackOrWhite() {
					proposedColor = nextProposedColor
					break
				}
			}
		}
	}

	return proposedColor
}

func minMax(n, min, max int) int {
	if n >= min && n <= max {
		return n
	} else if n < min {
		return min
	} else {
		return max
	}
}
