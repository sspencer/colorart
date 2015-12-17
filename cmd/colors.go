package main

//

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"time"

	"github.com/disintegration/gift"
	"github.com/sspencer/colorart"
)

type Cover struct {
	Filename, BackgroundColor, PrimaryColor, SecondaryColor, DetailColor string
}

func (c *Cover) String() string {
	return fmt.Sprintf("%s: bg=%s, primary=%s, secondary=%s, detail=%s",
		path.Base(c.Filename),
		c.BackgroundColor,
		c.PrimaryColor,
		c.SecondaryColor,
		c.DetailColor)
}

func analyzeFile(filename string, resize bool) (*Cover, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if resize {
		start := time.Now()
		g := gift.New(gift.Resize(500, 0, gift.LanczosResampling))
		dst := image.NewRGBA(g.Bounds(img.Bounds()))
		g.Draw(dst, img)
		img = dst
		fmt.Printf("- RESIZE %s took %s\n", path.Base(filename), time.Since(start))
	}

	start := time.Now()
	bg, c1, c2, c3 := colorart.Analyze(img)
	fmt.Printf("- ANALYZE %s took %s\n", path.Base(filename), time.Since(start))

	return &Cover{filename, bg.String(), c1.String(), c2.String(), c3.String()}, nil
}

func main() {

	if len(os.Args) < 2 {
		log.Fatalf("%s <img 1> <img 2> ... <img n>\n", os.Args[0])
	}

	/*
		f, err := os.Create("./cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	*/

	for i := 1; i < len(os.Args); i++ {
		cover, err := analyzeFile(os.Args[i], false)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(cover)
	}
}
