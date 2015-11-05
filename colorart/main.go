package main

//

import (
	"html/template"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"

	"github.com/sspencer/colorart"
)

type Cover struct {
	Filename, BackgroundColor, PrimaryColor, SecondaryColor, DetailColor string
}

func analyzeFile(filename string) (bg, c1, c2, c3 colorart.Color) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(os.Stderr, "%s: %v\n", "./cover.jpg", err)
	}

	return colorart.Analyze(img)
}

func main() {

	if len(os.Args) < 3 {
		log.Fatalf("%s <template> <img 1> <img 2> ... <img n>\n", os.Args[0])
	}

	tpl, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("webpage").Parse(string(tpl))

	covers := make([]Cover, 0, len(os.Args)-2)

	for i := 2; i < len(os.Args); i++ {
		bg, c1, c2, c3 := analyzeFile(os.Args[i])
		covers = append(covers, Cover{os.Args[i], bg.String(), c1.String(), c2.String(), c3.String()})
	}

	err = t.Execute(os.Stdout, covers)
	if err != nil {
		log.Fatal(err)
	}
}
