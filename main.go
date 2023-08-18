package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	setBottomRow()
	doFire()

	for i := 0; i < len(firePixels); i++ {
		termbox.SetCell(int(i%w), int(i/w), ' ',
			termbox.ColorDefault,
			termbox.Attribute(colorRamp[firePixels[i]]))
	}
	termbox.Flush()
}

func setBottomRow() {
	for x := 0; x < w; x++ {
		firePixels[h*w-x-1] = 11
	}
}

func doFire() {
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			spreadFire(y*w + x)
		}
	}
}

var w, h int
var firePixels []int
var colorRamp = [12]int{
	0,
	197,
	203,
	209,
	215,
	221,
	227,
	228,
	230,
	231,
	232,
	16, // White
}

func spreadFire(from int) {
	randomFactor := int(math.Round(rand.Float64()*3)) & 3
	to := from - w - randomFactor + 1
	if to < 0 {
		return
	}
	if to > (w*h)-1 {
		return
	}
	firePixels[to] = int(math.Max(0, float64((firePixels[from]-1)-(randomFactor&1))))
}

func main() {
	err := termbox.Init()

	w, h = termbox.Size()
	firePixels = make([]int, w*h)

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

	draw()

loop:
	for {
		select {
		case ev := <-event_queue:
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
				break loop
			}

		default:
			draw()
			time.Sleep(10 * time.Millisecond)
		}
	}
}
