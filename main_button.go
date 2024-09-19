package main

import (
    "image/color"
    "machine"
    "time"
)

type ColorWriter interface {
    WriteColors([]color.RGBA) error
    getLightness(int) float64
}

const (
    numPixels          = 20
    longPressThreshold = 500 * time.Millisecond
)

// var colorWriter ColorWriter - defined in main_tiny.go and main_biggy.go

type mode int

const (
    colorMode mode = iota
    speedMode
    // patternMode
    // lightnessMode
    modeLen
)

var modeLevels map[mode]int
var modeMaxLevels = map[mode]int{
    colorMode: 2,
    speedMode: 2,
}
var currentMode mode = colorMode
var minTick time.Duration = 50 * time.Millisecond

type ledConfig struct {
    mode  mode
    level int
    intro bool
}

func main() {
    writeColors(0, 0)
    pressCh := make(chan bool, 32)
    go startButtonListener(pressCh, machine.GP19)

    var hue int32 = 0
    hueCh := make(chan int32, 32)
    tickCh := make(chan time.Duration, 32)
    go startLEDs(hueCh, tickCh, machine.GP28)

    var pace time.Duration = 1
    for longPress := range pressCh {
        if longPress {
            pace = (pace + 1) % 5
            tickCh <- minTick * (pace + 1)
        } else {
            hue = (hue + 600) % 3600
            hueCh <- hue
        }
        // time.Sleep(200 * time.Millisecond)
        // colorsOff()
        // intro := false
        // if longPress {
        //     currentMode = (currentMode + 1) % modeLen // color, speed,...
        //     intro = true
        // } else { // short press
        //     modeLevels[currentMode] = (modeLevels[currentMode] + 1) % modeMaxLevels[currentMode]
        // }
        // ledConfigCh <- ledConfig{
        //     mode:  currentMode,
        //     level: modeLevels[currentMode],
        //     intro: intro,
        // }
    }
    select {}
}

func startLEDs(hueCh chan int32, tickCh chan time.Duration, _ machine.Pin) {
    var hue int32 = 1800
    tick := minTick
    idx := 0
    inc := 1
    //colorMultiplier := int32(3600 / modeMaxLevels[colorMode])
    for {
        select {
        case hue = <-hueCh:
        case tick = <-tickCh:
        case <-time.After(tick):
            writeColors(hue, idx)
            idx += inc
            if idx == 0 || idx == numPixels-1 {
                inc = -inc
            }
        }
    }
}

var colors = make([]color.RGBA, numPixels)

func writeColors(hue int32, idx int) error {
    for i := range colors {
        dist := i - idx
        if dist < 0 {
            dist = -dist
        }
        l := getLightness(dist)
        l = colorWriter.adjustLightness(l)
        colors[i] = hsl(hue, 1000, l)
    }
    return colorWriter.WriteColors(colors)
}

func colorsOff() error {
    for i := range colors {
        colors[i] = color.RGBA{}
    }
    return colorWriter.WriteColors(colors)

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
func startButtonListener(pressCh chan bool, btn machine.Pin) {
    start := time.Now()
    pressed := false
    ch := setupButtonPressChan(btn)
    for press := range ch { // button down (true) or up (false) event
        if press == pressed {
            continue
        }
        if press {
            start = time.Now()
        } else {
            pressCh <- time.Since(start) > longPressThreshold // long/short press
        }
        pressed = !pressed
    }
}

func setupButtonPressChan(btn machine.Pin) chan bool {
    config := machine.PinConfig{Mode: machine.PinInputPullup}
    ch := make(chan bool, 32)
    btn.Configure(config)
    btn.SetInterrupt(machine.PinFalling|machine.PinRising, func(pin machine.Pin) {
        select {
        case ch <- !pin.Get():
        default:
        }
    })
    return ch
}
