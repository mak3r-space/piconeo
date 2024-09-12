package main

import (
    "image/color"
    "math"
    "time"
)

type ColorWriter interface {
    WriteColors([]color.RGBA) error
    getLightness(int) float64
}

const numPixels = 20

// var colorWriter ColorWriter - defined in main_tiny.go and main_biggy.go

func main() {
    idx := 0
    inc := 1
    for true {
        if err := writeColors(idx, colorWriter); err != nil {
            println("ERROR: " + err.Error())
            return
        }
        time.Sleep(time.Millisecond * 100)
        idx += inc
        if idx == 0 || idx == numPixels-1 {
            inc = -inc
        }
    }
}

// hsl converts and HSL color to a color.RGBA.
// h is the hue, an angle in [0,360] s,l in [0,1]
// see: https://stackoverflow.com/a/64090995/661500
func hsl(h, s, l float64) color.RGBA {
    a := s * min(l, 1-l)
    f := func(n float64) uint8 {
        k := math.Mod(n+h/30, 12)
        v := l - a*max(min(k-3, 9-k, 1), -1)
        return uint8(v * 255)
    }
    return color.RGBA{R: f(0), G: f(8), B: f(4)}
}

func writeColors(idx int, cw ColorWriter) error {
    var colors = make([]color.RGBA, numPixels)
    fidx := float64(idx)
    for i := range colors {
        fi := float64(i)
        distIdx := int(math.Abs(fi - fidx))
        hue := fi / numPixels * 360
        lightness := cw.getLightness(distIdx)
        colors[i] = hsl(hue, 1, lightness)
    }
    return cw.WriteColors(colors)
}
