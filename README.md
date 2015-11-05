# ColorArt

Go port of [Panic's iTunes 11](https://www.panic.com/blog/itunes-11-and-colors/) album art [color algorithm](https://github.com/panicinc/ColorArt).

Looks like iTunes 12 uses a blurred version of the album cover now.  Here's a quick repro:

    gm mogrify -size 320x320 -format blur.jpg -blur 240x240 album.jpg

To view a demo:

    $ go run main.go covers.html ~/album/*.jpg > index.html

To do:

    * scale image before calculating colors
    * integrate into (my fork of) [imaginary](https://github.com/sspencer/imaginary).
