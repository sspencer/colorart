package colorart

import (
	"fmt"
	"image"
	"time"
)

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

	imageColors := NewCountedSet(2000)
	its := 0
	start := time.Now()
	for y := c.img.imgBounds.Min.Y; y < c.img.imgBounds.Max.Y; y++ {
		for x := c.img.imgBounds.Min.X; x < c.img.imgBounds.Max.X; x++ {
			imageColors.AddPixel(c.img.getPixel(x, y))
			its++
		}
	}
	fmt.Println("Iterations(1):", its, time.Since(start))
	findDarkTextColor := !backgroundColor.IsDarkColor()

	selectColors := NewCountedSet(1000)
	its = 0
	start = time.Now()

	for key, cnt := range imageColors.m {
		// don't bother unless there's more than a few of the same color
		if cnt > 10 {
			curColor := rgbToColor(key).ColorWithMinimumSaturation(0.15)
			if curColor.IsDarkColor() == findDarkTextColor {
				its++
				selectColors.AddCount(key, imageColors.Count(key))
				//fmt.Println(imageColors.Count(key))
			}
		}
	}
	fmt.Println("Iterations(2):", its, time.Since(start))

	sortedColors := selectColors.SortedSet()

	its = 0
	start = time.Now()
	for _, e := range sortedColors {
		its++
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
	fmt.Println("Iterations(3):", its, time.Since(start))

	return
}

func (c *colorArt) findEdgeColor() Color {

	edgeColors := NewCountedSet(500)
	x0 := c.img.imgBounds.Min.X
	x1 := c.img.imgBounds.Max.X - 1
	for y := c.img.imgBounds.Min.Y; y < c.img.imgBounds.Max.Y; y++ {
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
