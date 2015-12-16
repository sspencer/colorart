# ColorArt

Go port of [Panic's iTunes 11](https://www.panic.com/blog/itunes-11-and-colors/) album art [color algorithm](https://github.com/panicinc/ColorArt).

Looks like iTunes 12 uses a blurred version of the album cover now.  Here's a quick repro:

    gm mogrify -size 320x320 -format blur.jpg -blur 240x240 album.jpg

To view a demo:

    $ go run main.go covers.html ~/album/*.jpg > index.html

To speed things up, this code makes use of [GIFT](https://github.com/disintegration/gift) to resize images.  Also, the file "pixel.go"
from that project was copied directly into the project to make getting
pixels faster.

    go get -u github.com/disintegration/gift

