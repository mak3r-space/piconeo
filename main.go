//go:build !tinygo

package main

import (
    "image/color"
    "math"
    "time"
)

type ColorWriter interface {
    WriteColors([]color.RGBA) error
    adjustLightness(int32) int32
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

func writeColors(idx int, cw ColorWriter) error {
    var colors = make([]color.RGBA, numPixels)
    for i := range colors {
        dist := i - idx
        if dist < 0 {
            dist = -dist
        }
        h := int32(i) * 3600 / numPixels
        l := getLightness(dist)
        l = cw.adjustLightness(l)
        colors[i] = hsl(h, 1000, l)
    }
    return cw.WriteColors(colors)
}

func getLightness(dist int) int32 {
    lightness := []int32{500, 350, 200, 100, 50}
    if dist < len(lightness) {
        return lightness[dist]
    }
    return 0
}

// hsl converts a 10 x scaled HSL (Hue, Saturation, Lightness) color to an RGBA
// (Red, Green, Blue, Alpha) color.
//
// - h: The hue, represented as an integer in the range [0, 3600] (the angle*10 on the color wheel).
// - s: The saturation, represented as an integer in the range [0, 1000].
// - l: The lightness, represented as an integer in the range [0, 1000].
//
// Due to the use of integer calculations for efficiency, there might be
// occasional rounding errors of ±1 in the RGB components compared to
// floating-point implementations. For a reference implementation using
// floats, see the [hslFloat64] function below implemented according to
// https://stackoverflow.com/a/64090995/661500
func hsl(h, s, l int32) color.RGBA {
    a := s * min(l, 1000-l) / 1000
    f := func(n int32) uint8 {
        k := (n*300 + h) % 3600
        v := l - a*max(min(k-900, 2700-k, 300), -300)/300
        return uint8(v * 255 / 1000)
    }
    return color.RGBA{R: f(0), G: f(8), B: f(4)}
}

// hslFloat64 converts and HSL color to a color.RGBA.
// h is the hue, an angle in [0,360] s,l in [0,1]
// see: https://stackoverflow.com/a/64090995/661500
func hslFloat64(h, s, l float64) color.RGBA {
    a := s * min(l, 1-l)
    f := func(n float64) uint8 {
        k := math.Mod(n+h/30, 12)
        v := l - a*max(min(k-3, 9-k, 1), -1)
        return uint8(v * 255)
    }
    return color.RGBA{R: f(0), G: f(8), B: f(4)}
}
