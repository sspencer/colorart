package main

//

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/disintegration/gift"
	"github.com/sspencer/colorart"
)

const (
	doResizeImage   = true
	resizeThreshold = 210
	resizeSize      = 200
	tmpl            = `
<!DOCTYPE html>
<html>
<head>
	<title>Cover Art</title>
</head>
<body>
{{range .}}
<div style="background:{{.BackgroundColor}};height:240px;margin:4px;padding:4px">
	<img width="240" height="240" style="float:right" src="file://{{.Filename}}">
	<h1 style="color:{{.PrimaryColor}}">Primary Color</h1>
	<h2 style="color:{{.SecondaryColor}}">Secondary Color</h2>
	<h3 style="color:{{.DetailColor}}">Detail Color</h3>
</div>
{{end}}
</body>
</html>
`
)

type cover struct {
	Filename        string
	BackgroundColor colorart.Color
	PrimaryColor    colorart.Color
	SecondaryColor  colorart.Color
	DetailColor     colorart.Color
}

func (c *cover) String() string {
	return fmt.Sprintf("%s: bg=%s, primary=%s, secondary=%s, detail=%s",
		path.Base(c.Filename),
		c.BackgroundColor,
		c.PrimaryColor,
		c.SecondaryColor,
		c.DetailColor)
}

func analyzeFile(filename string) cover {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()
	if doResizeImage && (b.Max.X-b.Min.X > resizeThreshold || b.Max.Y-b.Min.Y > resizeThreshold) {
		g := gift.New(gift.Resize(resizeSize, 0, gift.LanczosResampling))
		dst := image.NewRGBA(image.Rect(0, 0, resizeSize, resizeSize))
		g.Draw(dst, img)
		img = dst
	}

	c := colorart.Analyze(img)

	return cover{filename, c.BackgroundColor, c.PrimaryColor, c.SecondaryColor, c.DetailColor}
}

func generateHtml(args []string) {
	t, err := template.New("webpage").Parse(string(tmpl))
	if err != nil {
		log.Fatalf("Error parsing template")
	}

	covers := make([]cover, 0, len(args))

	for _, arg := range args {
		cover := analyzeFile(arg)
		covers = append(covers, cover)
	}

	err = t.Execute(os.Stdout, covers)
	if err != nil {
		log.Fatal(err)
	}
}

func generateText(args []string) {
	for _, arg := range args {
		cover := analyzeFile(arg)
		fmt.Println(cover)
	}
}

func main() {

	htmlPtr := flag.Bool("html", false, "html output")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatalf("%s [-html] <img 1> <img 2> ... <img n>\n", os.Args[0])
	}

	if *htmlPtr == true {
		generateHtml(args)
	} else {
		generateText(args)
	}
}
