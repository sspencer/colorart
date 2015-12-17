package colorart

import "image"

type colorArt struct {
	img *pixelGetter
}

// resize image
// https://github.com/nfnt/resize

func Analyze(img image.Image) (backgroundColor, primaryColor, secondaryColor, detailColor Color) {
	c := &colorArt{}
	c.img = newPixelGetter(img)

	backgroundColor = c.findEdgeColor()
	primaryColor, secondaryColor, detailColor = c.findTextColors(backgroundColor)

	darkBackground := backgroundColor.IsDarkColor()

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

	return
}

func (c *colorArt) findTextColors(backgroundColor Color) (primaryColor, secondaryColor, detailColor Color) {
	b := c.img.imgBounds
	imageColors := parallelize(b.Min.Y, b.Max.Y, func(ch chan *CountedSet, pmin, pmax int) {
		b := c.img.imgBounds
		colors := NewCountedSet(10000)
		for y := pmin; y < pmax; y += 1 {
			for x := b.Min.X; x < b.Max.X; x += 1 {
				colors.AddPixel(c.img.getPixel(x, y))
			}
		}

		ch <- colors
	})

	useDarkTextColor := !backgroundColor.IsDarkColor()
	selectColors := NewCountedSet(5000)

	for key, cnt := range imageColors.m {
		// don't bother unless there's more than a few of the same color
		if cnt > 10 {
			curColor := rgbToColor(key).ColorWithMinimumSaturation(0.15)
			if curColor.IsDarkColor() == useDarkTextColor {
				selectColors.AddCount(key, cnt)
			}
		}
	}

	sortedColors := selectColors.SortedSet()
	for _, e := range sortedColors {
		curColor := rgbToColor(e.Color)
		if !primaryColor.set {
			if curColor.IsContrastingColor(backgroundColor) {
				primaryColor = curColor
			}
		} else if !secondaryColor.set {
			if !primaryColor.IsDistinctColor(curColor) || !curColor.IsContrastingColor(backgroundColor) {
				continue
			}
			secondaryColor = curColor

		} else if !detailColor.set {
			if !secondaryColor.IsDistinctColor(curColor) ||
				!primaryColor.IsDistinctColor(curColor) ||
				!curColor.IsContrastingColor(backgroundColor) {
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
	if proposedColor.IsBlackOrWhite() {
		for i, e := range sortedColors {
			if i == 0 {
				// first entry is already set as "proposedEntry"
				continue
			}

			nextProposedEntry := e
			// make sure second choice is 30% as common as first choice
			if float32(nextProposedEntry.Count)/float32(proposedEntry.Count) > 0.3 {
				nextProposedColor := rgbToColor(nextProposedEntry.Color)
				if !nextProposedColor.IsBlackOrWhite() {
					proposedColor = nextProposedColor
					break
				}
			}
		}
	}

	return proposedColor
}
