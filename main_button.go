package main

import (
    "image/color"
    "machine"
    "time"
)

type ColorWriter interface {
    WriteColors([]color.RGBA) error
    adjustLightness(int32) int32
}

const (
    numPixels          = 20
    longPressThreshold = 500 * time.Millisecond
)

// var colorWriter ColorWriter - defined in main_tiny.go and main_biggy.go

type mode int

const (
    paletteMode mode = iota
    speedMode
    // patternMode
    lightnessMode
    modeLen
)

var levels = make(map[mode]int, modeLen)
var lightnessLevels = [][]int32{
    {500, 350, 200, 100, 50},
    {200, 140, 80, 40, 20},
    {100, 70, 40, 25, 10},
    {70, 40, 25, 10, 5},
}
var palettes = [][]int32{{0, 3600}, {0, 1200}, {1200, 2400}, {2400, 3600}}
var levelCnt = map[mode]int{
    paletteMode:   len(palettes),
    lightnessMode: len(lightnessLevels),
    speedMode:     4,
}
var cur mode = paletteMode
var minTick time.Duration = 75 * time.Millisecond

type ledConfig struct {
    mode      mode
    level     int
    modeIntro bool
}

func main() {
    writeColors(0, 0, 0, 0)
    modeCh := make(chan bool, 32)
    go startButtonListener(modeCh, machine.GP16)
    levelCh := make(chan bool, 32)
    go startButtonListener(levelCh, machine.GP26)
    configCh := make(chan ledConfig, 32)
    go startLEDs(configCh, machine.GP28)

    for {
        select {
        case <-modeCh:
            cur = (cur + 1) % modeLen
        case <-levelCh:
            levels[cur] = (levels[cur] + 1) % levelCnt[cur]
            configCh <- ledConfig{mode: cur, level: levels[cur]}
        }
    }
}

func startLEDs(configCh chan ledConfig, _ machine.Pin) {
    var paletteStart, paletteEnd int32
    lightness := 0
    tick := minTick
    idx := 0
    inc := 1
    for {
        select {
        case config := <-configCh:
            switch config.mode {
            case paletteMode:
                paletteStart = palettes[config.level][0]
                paletteEnd = palettes[config.level][1]
            case speedMode:
                tick = minTick * time.Duration(1<<config.level)
            case lightnessMode:
                lightness = config.level
            }
        case <-time.After(tick):
            writeColors(paletteStart, paletteEnd, lightness, idx)
            idx += inc
            if idx == 0 || idx == numPixels-1 {
                inc = -inc
            }
        }
    }
}

var colors = make([]color.RGBA, numPixels)

func writeColors(paletteStart int32, paletteEnd int32, lightness int, idx int) error {
    for i := range colors {
        dist := i - idx
        if dist < 0 {
            dist = -dist
        }
        l := getLightness(dist, lightness)
        l = colorWriter.adjustLightness(l)
        hue := getHue(paletteStart, paletteEnd, int32(i))
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

func getHue(paletteStart int32, paletteEnd int32, idx int32) int32 {
    return paletteStart + (paletteEnd-paletteStart)*idx/numPixels
}

func getLightness(dist int, level int) int32 {
    lightness := lightnessLevels[level]
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
