package main

//

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"

	"github.com/disintegration/gift"
)

const (
	resizeThreshold = 350
	resizeSize      = 320
	sigma           = 40
)

func doit(fn string) (string, error) {
	file, err := os.Open(fn)

	if err != nil {
		return "", err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	b := img.Bounds()
	var g *gift.GIFT

	if b.Max.X-b.Min.X >= resizeThreshold || b.Max.Y-b.Min.Y >= resizeThreshold {
		g = gift.New(
			gift.Resize(resizeSize, resizeSize, gift.LanczosResampling),
			gift.GaussianBlur(sigma))
	} else {
		g = gift.New(gift.GaussianBlur(sigma))
	}

	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)
	img = dst

	fn = path.Base(fn)
	ext := path.Ext(fn)
	fn = fmt.Sprintf("./%s.blur%s", fn[:len(fn)-len(ext)], ext)

	w, _ := os.Create(fn)
	defer w.Close()
	if err = jpeg.Encode(w, dst, nil); err != nil {
		return "", err
	}

	return fn, nil
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("%s <img 1> <img 2> ... <img n>\n", os.Args[0])
	}

	for i := 1; i < len(os.Args); i++ {
		fn, err := doit(os.Args[i])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Wrote file:", fn)
	}
}
