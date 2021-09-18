//hershey fonts drawing tests

package hershey

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/StephaneBunel/bresenham"
)

var ColorWHITE = color.RGBA{255, 255, 255, 255}
var ColorBLACK = color.RGBA{0, 0, 0, 255}
var ColorGREEN = color.RGBA{0, 255, 0, 255}

//Test function for DrawChar
func DrawAllFontImage() (err error) {
	var imgRect = image.Rect(0, 0, 2000, 2000)
	img := image.NewRGBA(imgRect)
	var mx, my int //reserve space for move coords
	x := 50
	y := 1900
	lineHt := 50
	scale := 1
	for fontname, f := range Fonts {
		fmt.Println("  Drawing Font: ", fontname)
		err = DrawString(fontname+":", "Plain", scale, &x, &y, ImageMoveTo, ImageLineTo, &mx, &my, img, &ColorBLACK)
		x = 50
		y -= Height["Plain"][1] - Height["Plain"][0] + 5
		for i, _ := range f {
			err = DrawChar(rune(i+32), fontname, scale, &x, &y, ImageMoveTo, ImageLineTo, &mx, &my, img, &ColorBLACK)
			if err != nil {
				return
			}
		}
		y -= lineHt
		x = 50
	}
	toimg, err := os.Create("allfonts.png")
	if err != nil {
		return
	}
	defer toimg.Close()
	flipImg := FlipV(img)
	png.Encode(toimg, flipImg)
	return
}

func DrawAllFontStringImage() (err error) {
	var imgRect = image.Rect(0, 0, 2000, 2000)
	img := image.NewRGBA(imgRect)
	var mx, my int //reserve space for move coords
	x := 0
	y := 1990
	scale := 2
	str := "Hello World!"
	for fontname, _ := range Fonts {
		fmt.Println("  Drawing Font String: ", fontname)
		//miny, maxy
		height := Height[fontname]
		y -= (height[1] * scale)
		x = 0
		err = DrawString(fontname+": "+str, fontname, scale, &x, &y, ImageMoveTo, ImageLineTo, &mx, &my, img, &ColorBLACK)
		if err != nil {
			return err
		}
		y += (height[0] * scale)
	}
	toimg, err := os.Create("allstringfonts.png")
	if err != nil {
		return
	}
	defer toimg.Close()
	flipImg := FlipV(img)
	png.Encode(toimg, flipImg)
	return
}

//Test function for DrawStringLinesImage
func DrawAStringLines() (err error) {
	var imgRect = image.Rect(0, 0, 3000, 3000)
	img := image.NewRGBA(imgRect)
	var mx, my int //reserve space for move coords
	x := 100
	y := 2900
	scale := 3
	lineSpace := 50
	width := 2800
	font := "Script_Complex"
	str := "Dearest Bianca,\n\nHello favorite cousin, how are you doing? I am writing to you from sunny Seattle where the trains and ferries are always running and the seagulls are always calling. We are having loads of fun. I have decided that I really want to be submarine captain for my career, so I am now studying Ocean Navigation. Perhaps you would like to travel with me when I get my submarine?\nLove,\n Gopher"
	err = DrawStringLines(str, font, scale, &x, &y, lineSpace, width, ImageMoveTo, ImageLineTo, &mx, &my, img, &ColorBLACK)
	if err != nil {
		return err
	}
	toimg, err := os.Create("stringlines.png")
	if err != nil {
		return
	}
	defer toimg.Close()
	flipImg := FlipV(img)
	png.Encode(toimg, flipImg)
	return
}

//find better way to convert interface?
func convertIf(img image.Image) draw.Image {
	return img.(draw.Image)
}

// s = (x2, y2, x1, y1, img, color)
func ImageMoveTo(s ...interface{}) {
	x := *(s[0].(*int))
	y := *(s[1].(*int))
	*(s[2].(*int)) = x //save x and y in s for line
	*(s[3].(*int)) = y
}

// s = (x2, y2, x1, y1, img, color)
func ImageLineTo(s ...interface{}) {
	x1 := *(s[2].(*int)) //saved from ImageMoveTo
	y1 := *(s[3].(*int))
	x2 := *(s[0].(*int))
	y2 := *(s[1].(*int))
	img := (s[4].(*image.RGBA))
	dimg := convertIf(img) //bresenham needs draw.Image
	color := (s[5].(*color.RGBA))
	bresenham.Bresenham(dimg, x1, y1, x2, y2, color)
	*(s[2].(*int)) = x2 //save x and y in s for next line
	*(s[3].(*int)) = y2
}

func FlipV(img image.Image) *image.RGBA {
	bnds := img.Bounds()
	var newImg = image.NewRGBA(bnds)
	for j := bnds.Min.Y; j < bnds.Max.Y; j++ {
		for i := bnds.Min.X; i < bnds.Max.X; i++ {
			c := img.At(i, j)
			newImg.Set(i, bnds.Max.Y-j-1, c)
		}
	}
	return newImg
}
