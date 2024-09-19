//go:build !tinygo

package main

import (
    "fmt"
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

func main2() {
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

func main() {
    testHSL("red", 0, 1, 0.5)
    testHSL("green", 120, 1, 0.5)
    testHSL("magenta", 300, 1, 0.5)
    testHSL("light cyan", 180, 1, 0.8)
    testHSL("grey", 300, 0, 0.8)
}

func testHSL(label string, h, s, l float64) {
    fmt.Printf("\n%s:\n", label)
    rgb := hsl(h, s, l)
    fmt.Printf(" hsl(%7.0f %7.2f %7.2f) as RGBA %v\n", h, s, l, rgb)
    h2 := int(h * 10)
    s2 := int(s * 1000)
    l2 := int(l * 1000)
    rgb2 := hsl2(h2, s2, l2)
    // s = s * 1000
    fmt.Printf("hsl2(%7d %7d %7d) as RGBA %v\n", h2, s2, l2, rgb2)
    if rgb != rgb2 {
        fmt.Printf("💥 %v != %v\n", rgb, rgb2)
    }
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

// hsl2 converts and HSL color to a color.RGBA.
// h is the hue, an angle in [0,3600] s,l in [0,1000]
// see: https://stackoverflow.com/a/64090995/661500
func hsl2(h, s, l int) color.RGBA {
    l2 := float64(l)
    a2 := float64(s) * min(l2, 1000-l2) / 1000

    a := s * min(l, 1000-l) / 1000
    f := func(n int) uint8 {
        //k := math.Mod(n+h/3, 1200) / 100
        //k := math.Mod(n+h, 3600) / 300
        k := (n + h) % 3600 / 300
        k2 := float64(k)
        _ = (l2 - a2*max(min(k2-3, 9-k2, 1), -1))

        v := l - a*max(min(k-3, 9-k, 1), -1)

        return uint8(v * 255 / 1000)
    }
    return color.RGBA{R: f(0), G: f(2400), B: f(1200)}
}
