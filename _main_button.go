package main

import (
    "image/color"
    "machine"
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
    green := 120.0
    red := 0.0
    purple := 300.0
    writeColors(purple, colorWriter) // blue

    ////
    pressCh := make(chan bool, 32)
    go buttonListener(pressCh, machine.GP19)
    for longPress := range pressCh {
        if longPress {
            writeColors(green, colorWriter) // long press
            mode = (mode + 1) % 3           // pattern, color, birghtness
        } else { // short press
            level[mode] = (level[mode] + 1) % 4
            updateLEDConfig(mode, level[mode])
            //writeColors(red, colorWriter)
        }
        // time.Sleep(time.Millisecond * 100)
        //writeOff(colorWriter)
    }

    select {}
}

func buttonListener(pressCh chan bool, btn machine.Pin) {
    start := time.Now()
    pressed := false
    config := machine.PinConfig{Mode: machine.PinInputPullup}
    ch := make(chan bool, 32)
    btn.Configure(config)
    btn.SetInterrupt(machine.PinFalling|machine.PinRising, func(pin machine.Pin) {
        select {
        case ch <- !pin.Get():
        default:
        }
    })

    for press := range ch {
        if press == pressed {
            continue
        }
        if press {
            start = time.Now()
        } else {
            pressCh <- time.Since(start) > 500*time.Millisecond // long/short press
        }
        pressed = !pressed
    }
}

//func buttonTest(cw ColorWriter) {
// if !pin.Get() {
//     start = time.Now()
//     writeColors(math.Mod(float64(start.Second()), 360.0), cw)
//     return
// }
// dur := time.Since(start)
// long := false
// if dur > time.Millisecond*500 {
//     long = true
// }
// led.Set(true)
// time.Sleep(time.Millisecond * 200)
// if long {
//     led.Set(false)
//     time.Sleep(time.Millisecond * 200)
//     led.Set(true)
//     time.Sleep(time.Millisecond * 200)
// }

//original
// writeColors(float64(color), colorWriter)
// color = (color + 60) % 360
// led.Set(!pin.Get())
// })

//}

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

func writeColors(hue float64, cw ColorWriter) error {
    var colors = make([]color.RGBA, numPixels)
    for i := range colors {
        colors[i] = hsl(hue, 1, 0.5)
    }
    return cw.WriteColors(colors)
}

func writeOff(cw ColorWriter) error {
    var colors = make([]color.RGBA, numPixels)
    for i := range colors {
        colors[i] = hsl(0, 1, 0)
    }
    return cw.WriteColors(colors)
}

func writeColorsOrig(idx int, cw ColorWriter) error {
    var colors = make([]color.RGBA, numPixels)
    fidx := float64(idx)
    for i := range colors {
        fi := float64(i)
        distIdx := int(math.Abs(fi - fidx))
        hue := 0.0 // fi / numPixels * 360
        lightness := cw.getLightness(distIdx)
        colors[i] = hsl(hue, 1, lightness)
    }
    return cw.WriteColors(colors)
}
