/*
The MIT License (MIT)

Copyright (c) 2014 Grigory Dryapak

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package colorart

import (
	"image"
	"image/color"
)

// code is from github.com/disintegration/gift/pixels.go
type pixel struct {
	R, G, B, A float32
}

type imageType int

const (
	itGeneric imageType = iota
	itNRGBA
	itNRGBA64
	itRGBA
	itRGBA64
	itYCbCr
	itGray
	itGray16
	itPaletted
)

type pixelGetter struct {
	imgType     imageType
	imgBounds   image.Rectangle
	imgGeneric  image.Image
	imgNRGBA    *image.NRGBA
	imgNRGBA64  *image.NRGBA64
	imgRGBA     *image.RGBA
	imgRGBA64   *image.RGBA64
	imgYCbCr    *image.YCbCr
	imgGray     *image.Gray
	imgGray16   *image.Gray16
	imgPaletted *image.Paletted
	imgPalette  []pixel
}

func newPixelGetter(img image.Image) (p *pixelGetter) {
	switch img := img.(type) {
	case *image.NRGBA:
		p = &pixelGetter{
			imgType:   itNRGBA,
			imgBounds: img.Bounds(),
			imgNRGBA:  img,
		}

	case *image.NRGBA64:
		p = &pixelGetter{
			imgType:    itNRGBA64,
			imgBounds:  img.Bounds(),
			imgNRGBA64: img,
		}

	case *image.RGBA:
		p = &pixelGetter{
			imgType:   itRGBA,
			imgBounds: img.Bounds(),
			imgRGBA:   img,
		}

	case *image.RGBA64:
		p = &pixelGetter{
			imgType:   itRGBA64,
			imgBounds: img.Bounds(),
			imgRGBA64: img,
		}

	case *image.Gray:
		p = &pixelGetter{
			imgType:   itGray,
			imgBounds: img.Bounds(),
			imgGray:   img,
		}

	case *image.Gray16:
		p = &pixelGetter{
			imgType:   itGray16,
			imgBounds: img.Bounds(),
			imgGray16: img,
		}

	case *image.YCbCr:
		p = &pixelGetter{
			imgType:   itYCbCr,
			imgBounds: img.Bounds(),
			imgYCbCr:  img,
		}

	case *image.Paletted:
		p = &pixelGetter{
			imgType:     itPaletted,
			imgBounds:   img.Bounds(),
			imgPaletted: img,
			imgPalette:  convertPalette(img.Palette),
		}
		return

	default:
		p = &pixelGetter{
			imgType:    itGeneric,
			imgBounds:  img.Bounds(),
			imgGeneric: img,
		}
	}
	return
}

const (
	qf8  = float32(1.0 / 255.0)
	qf16 = float32(1.0 / 65535.0)
	epal = qf16 * qf16 / 2.0
)

func convertPalette(p []color.Color) []pixel {
	plen := len(p)
	pnew := make([]pixel, plen)
	for i := 0; i < plen; i++ {
		r16, g16, b16, a16 := p[i].RGBA()
		switch a16 {
		case 0:
			pnew[i] = pixel{0.0, 0.0, 0.0, 0.0}
		case 65535:
			r := float32(r16) * qf16
			g := float32(g16) * qf16
			b := float32(b16) * qf16
			pnew[i] = pixel{r, g, b, 1.0}
		default:
			q := float32(1.0) / float32(a16)
			r := float32(r16) * q
			g := float32(g16) * q
			b := float32(b16) * q
			a := float32(a16) * qf16
			pnew[i] = pixel{r, g, b, a}
		}
	}
	return pnew
}

func pixelclr(c color.Color) (px pixel) {
	r16, g16, b16, a16 := c.RGBA()
	switch a16 {
	case 0:
		px = pixel{0.0, 0.0, 0.0, 0.0}
	case 65535:
		r := float32(r16) * qf16
		g := float32(g16) * qf16
		b := float32(b16) * qf16
		px = pixel{r, g, b, 1.0}
	default:
		q := float32(1.0) / float32(a16)
		r := float32(r16) * q
		g := float32(g16) * q
		b := float32(b16) * q
		a := float32(a16) * qf16
		px = pixel{r, g, b, a}
	}
	return px
}

func (p *pixelGetter) getPixel(x, y int) (px pixel) {
	switch p.imgType {
	case itNRGBA:
		i := p.imgNRGBA.PixOffset(x, y)
		r := float32(p.imgNRGBA.Pix[i+0]) * qf8
		g := float32(p.imgNRGBA.Pix[i+1]) * qf8
		b := float32(p.imgNRGBA.Pix[i+2]) * qf8
		a := float32(p.imgNRGBA.Pix[i+3]) * qf8
		px = pixel{r, g, b, a}

	case itNRGBA64:
		i := p.imgNRGBA64.PixOffset(x, y)
		r := float32(uint16(p.imgNRGBA64.Pix[i+0])<<8|uint16(p.imgNRGBA64.Pix[i+1])) * qf16
		g := float32(uint16(p.imgNRGBA64.Pix[i+2])<<8|uint16(p.imgNRGBA64.Pix[i+3])) * qf16
		b := float32(uint16(p.imgNRGBA64.Pix[i+4])<<8|uint16(p.imgNRGBA64.Pix[i+5])) * qf16
		a := float32(uint16(p.imgNRGBA64.Pix[i+6])<<8|uint16(p.imgNRGBA64.Pix[i+7])) * qf16
		px = pixel{r, g, b, a}

	case itRGBA:
		i := p.imgRGBA.PixOffset(x, y)
		a8 := p.imgRGBA.Pix[i+3]
		switch a8 {
		case 0:
			px = pixel{0.0, 0.0, 0.0, 0.0}
		case 255:
			r := float32(p.imgRGBA.Pix[i+0]) * qf8
			g := float32(p.imgRGBA.Pix[i+1]) * qf8
			b := float32(p.imgRGBA.Pix[i+2]) * qf8
			px = pixel{r, g, b, 1.0}
		default:
			q := float32(1.0) / float32(a8)
			r := float32(p.imgRGBA.Pix[i+0]) * q
			g := float32(p.imgRGBA.Pix[i+1]) * q
			b := float32(p.imgRGBA.Pix[i+2]) * q
			a := float32(a8) * qf8
			px = pixel{r, g, b, a}
		}

	case itRGBA64:
		i := p.imgRGBA64.PixOffset(x, y)
		a16 := uint16(p.imgRGBA64.Pix[i+6])<<8 | uint16(p.imgRGBA64.Pix[i+7])
		switch a16 {
		case 0:
			px = pixel{0.0, 0.0, 0.0, 0.0}
		case 65535:
			r := float32(uint16(p.imgRGBA64.Pix[i+0])<<8|uint16(p.imgRGBA64.Pix[i+1])) * qf16
			g := float32(uint16(p.imgRGBA64.Pix[i+2])<<8|uint16(p.imgRGBA64.Pix[i+3])) * qf16
			b := float32(uint16(p.imgRGBA64.Pix[i+4])<<8|uint16(p.imgRGBA64.Pix[i+5])) * qf16
			px = pixel{r, g, b, 1.0}
		default:
			q := float32(1.0) / float32(a16)
			r := float32(uint16(p.imgRGBA64.Pix[i+0])<<8|uint16(p.imgRGBA64.Pix[i+1])) * q
			g := float32(uint16(p.imgRGBA64.Pix[i+2])<<8|uint16(p.imgRGBA64.Pix[i+3])) * q
			b := float32(uint16(p.imgRGBA64.Pix[i+4])<<8|uint16(p.imgRGBA64.Pix[i+5])) * q
			a := float32(a16) * qf16
			px = pixel{r, g, b, a}
		}

	case itGray:
		i := p.imgGray.PixOffset(x, y)
		v := float32(p.imgGray.Pix[i]) * qf8
		px = pixel{v, v, v, 1.0}

	case itGray16:
		i := p.imgGray16.PixOffset(x, y)
		v := float32(uint16(p.imgGray16.Pix[i+0])<<8|uint16(p.imgGray16.Pix[i+1])) * qf16
		px = pixel{v, v, v, 1.0}

	case itYCbCr:
		iy := p.imgYCbCr.YOffset(x, y)
		ic := p.imgYCbCr.COffset(x, y)
		r8, g8, b8 := color.YCbCrToRGB(p.imgYCbCr.Y[iy], p.imgYCbCr.Cb[ic], p.imgYCbCr.Cr[ic])
		r := float32(r8) * qf8
		g := float32(g8) * qf8
		b := float32(b8) * qf8
		px = pixel{r, g, b, 1.0}

	case itPaletted:
		i := p.imgPaletted.PixOffset(x, y)
		k := p.imgPaletted.Pix[i]
		px = p.imgPalette[k]

	case itGeneric:
		px = pixelclr(p.imgGeneric.At(x, y))
	}
	return
}

func (p *pixelGetter) getPixelRow(y int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for x := p.imgBounds.Min.X; x != p.imgBounds.Max.X; x++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

func (p *pixelGetter) getPixelColumn(x int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for y := p.imgBounds.Min.Y; y != p.imgBounds.Max.Y; y++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}
