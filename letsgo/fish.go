package letsgo

import (
	"fmt"
	"image"
	"sync"

	"keybd_event"
	"math"
	"time"

	"github.com/vova616/screenshot"
)

type Fish struct {
	lock   *sync.Mutex
	config Configeration
	kb     keybd_event.KeyBonding
	times  int
	ch     *chan int
}

func NewFish(config Configeration, times int) *Fish {
	var lock *sync.Mutex = new(sync.Mutex)
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	kb.SetKeys(keybd_event.VK_Q)
	ch := make(chan int, 1)
	return &Fish{
		config: config,
		lock:   lock,
		kb:     kb,
		times:  times,
		ch:     &ch,
	}
}

func (this *Fish) search(img *image.RGBA, find *bool, sx, ex, sy, ey int) {
x:
	for y := sy; y < ey; y++ {
		for x := sx; x <= ex; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if this.rgbInterval(r, g, b) {
				this.putValue(find, x, y)
				break x
			}
		}
	}
}

func (this *Fish) putValue(find *bool, x, y int) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if !*find {
		*this.ch <- x
		*this.ch <- y
		*find = true
	}
}

func (this *Fish) Launch() {
	if this.times <= 0 {
		for {
			this.fish()
			time.Sleep(time.Duration(this.config.Interval) * time.Second)
		}
	} else {
		for i := 0; i < this.times; i++ {
			this.fish()
			time.Sleep(time.Duration(this.config.Interval) * time.Second)
		}
	}
}

func (this *Fish) fish() {
	fmt.Println("开始钓鱼...")
	var find *bool
	*find = false
	this.kb.Launching()
	time.Sleep(time.Duration(this.config.SpellTime) * time.Second)

	n := time.Now().UnixNano()
	for t := 0; t < 1; t++ {
		if *find {
			break
		}
		go func() {
			for !*find {
				rect := image.Rect(this.config.StartX, this.config.StartY, this.config.StartX+this.config.Round*2, this.config.StartY+this.config.Round*2)
				img, err := screenshot.CaptureRect(rect)
				if err != nil {
					panic(err)
				}
				bounds := img.Bounds()
				//左下
				go this.search(img, find, bounds.Min.X, bounds.Max.X-this.config.Round, bounds.Min.Y, bounds.Max.Y-this.config.Round)
				//右下
				go this.search(img, find, bounds.Min.X+this.config.Round, bounds.Max.X, bounds.Min.Y, bounds.Max.Y-this.config.Round)
				//左上
				go this.search(img, find, bounds.Min.X, bounds.Max.X-this.config.Round, bounds.Min.Y+this.config.Round, bounds.Max.Y)
				//右上
				go this.search(img, find, bounds.Min.X+this.config.Round, bounds.Max.X, bounds.Min.Y+this.config.Round, bounds.Max.Y)
				//3秒还未找到，直接跳过
				if time.Now().UnixNano()-n > 3000000000 {
					this.putValue(find, 0, 0)
				}
			}
		}()
	}
	rx := <-*this.ch
	ry := <-*this.ch

	fmt.Println("找到中线相对位置，", rx, ry)
	angle := this.getAngle(this.config.Round, this.config.Round, rx, ry)
	ratio := angle / 360.0
	pressWaiting := this.config.RoundTime * ratio
	pressWaiting = math.Ceil(pressWaiting * 1000 * 1000 * 1000)
	time.Sleep(time.Duration(pressWaiting) * time.Nanosecond)
	this.kb.Launching()
	fmt.Println("钓鱼完成...")
}

func (this *Fish) rgbInterval(r, g, b uint32) bool {
	return this.config.Rmin <= int(r>>8) && int(r>>8) <= this.config.Rmax && this.config.Gmin <= int(g>>8) && int(g>>8) <= this.config.Gmax && this.config.Bmin <= int(b>>8) && int(b>>8) <= this.config.Bmax
}

func (this *Fish) getAngle(x1, y1, x2, y2 int) float64 {
	angle := 0.0
	// 直角的边长
	var x = float64(x1 - x2)
	var y = float64(y1 - y2)
	bearingRadians := math.Atan2(y, x)
	bearingDegrees := bearingRadians * (180.0 / math.Pi)
	if bearingDegrees <= 0 {
		angle = 360.0 + bearingDegrees
	} else {
		angle = bearingDegrees
	}
	if angle >= 90 {
		angle = angle - 90
	} else {
		angle = angle + 270
	}
	return angle
}
