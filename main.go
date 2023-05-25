package main

import (
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
)

const (
	ascii = "@%#*+=-:. "
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("There was no image pass"))
	}
	ipath := os.Args[1]
	im, err := os.Open(ipath)
	if err != nil {
		log.Fatal(err)
	}
	defer im.Close()
	img, _, err := image.Decode(im)
	if err != nil {
		log.Fatal(err)
	}
	asciiImg := convertToascii(imageResize(img, 150))
	name := parseName(ipath)
	f, err := os.Create(name + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	f.Write([]byte(asciiImg))
	f.Close()
}
func parseName(imagePath string) string {
	sp := strings.Split(imagePath, ".")
	sp1 := strings.Split(sp[0], "/")
	return sp1[len(sp1)-1]
}
func convertToascii(img image.Image) string {
	bounds := img.Bounds()
	with := bounds.Dx()
	height := bounds.Dy()
	var asciiArt strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < with; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			idx := int(gray.Y) * (len(ascii) - 1) / 255
			asciiArt.WriteString(string(ascii[idx]))
		}
		asciiArt.WriteString("\n")
	}
	return asciiArt.String()
}
func imageResample(dest *image.RGBA, src image.Image) {
	bounds := src.Bounds()
	srcWith := bounds.Dx()
	srcHeight := bounds.Dy()
	destHeight := dest.Bounds().Dy()
	destWith := dest.Bounds().Dx()
	for y := 0; y < destHeight; y++ {
		for x := 0; x < destWith; x++ {
			srcX := x * srcWith / destHeight
			srcY := y * srcHeight / destHeight
			dest.Set(x, y, src.At(srcX, srcY))
		}
	}
}
func imageResize(img image.Image, with int) image.Image {
	aspectRatio := float64(img.Bounds().Dx()) / float64(img.Bounds().Dy())
	height := with / int(aspectRatio)
	resizedImg := image.NewRGBA(image.Rect(0, 0, with, height))
	imageResample(resizedImg, img)
	return resizedImg
}
