package letsgo

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

var rmin, rmax, gmin, gmax, bmin, bmax = 240, 250, 230, 240, 128, 138

const (
	r = 70
)

func main1() {
	startX := 580
	startY := 308
	roundTime := 2.0
	time.Sleep(10 * time.Second)

	robotgo.KeyTap("q")

	// r := image.Rect(50, 50, 500, 500)
	// img, err := screenshot.CaptureRect(r)

	f1, err := os.Open("./2.png")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	origin, _, err := image.Decode(f1)
	if err != nil {
		panic(err)
	}
	img := origin.(*image.RGBA)
	subImg := img.SubImage(image.Rect(startX, startY, startX+r*2, startY+r*2)).(*image.RGBA)
	bounds := subImg.Bounds()
	rx := 0
	ry := 0
	fmt.Println(bounds.Min.X, bounds.Max.X, bounds.Max.Y)
b:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x <= bounds.Max.X; x++ {
			r, g, b, _ := subImg.At(x, y).RGBA()
			if rgbInterval(r, g, b) {
				ry = y
				rx = x
				break b
			}
		}
	}
	rx = rx - bounds.Min.X
	ry = ry - bounds.Min.Y
	fmt.Println(rx, ry)
	if rx < 0 || ry < 0 {
		return
	}
	angle := getAngle(r, r, rx, ry)

	ratio := angle / 360.0
	fmt.Println(ratio)
	pressWaiting := roundTime * ratio
	pressWaiting = pressWaiting * 1000 * 1000 * 1000
	time.Sleep(time.Duration(pressWaiting) * time.Nanosecond)
	f, err := os.Create("./ss.png")

	if err != nil {
		panic(err)
	}
	err = png.Encode(f, subImg)

	if err != nil {
		panic(err)
	}
	f.Close()
}

func rgbInterval(r, g, b uint32) bool {
	return rmin <= int(r>>8) && int(r>>8) <= rmax && gmin <= int(g>>8) && int(g>>8) <= gmax && bmin <= int(b>>8) && int(b>>8) <= bmax
}

func getAngle(x1, y1, x2, y2 int) float64 {
	// 直角的边长
	var x = float64(x1 - x2)
	var y = float64(y1 - y2)
	bearingRadians := math.Atan2(y, x)
	bearingDegrees := bearingRadians * (180.0 / math.Pi)
	if bearingDegrees <= 0 {
		return 360.0 + bearingDegrees
	}
	return bearingDegrees - 90
}
