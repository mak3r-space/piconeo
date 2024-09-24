package main

import (
    "image/color"
    "machine"
    "time"
)

type ColorWriter interface {
    WriteColors([]color.RGBA) error
    adjustLightness(int) int
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
    lightnessMode
    modeLen
)

type palette struct{ start, end int }

type ledConfig struct {
    mode      mode
    level     int
    modeIntro bool
}
type hsl struct{ H, S, L int }

var levels = make(map[mode]int, modeLen)
var lightnessLevels = [][]int{
    {500, 350, 200, 100, 50},
    {200, 140, 80, 40, 20},
    {100, 70, 40, 25, 10},
    {70, 40, 25, 10, 5},
}
var palettes = []palette{
    {start: 0, end: 3600},   // rainbow
    {start: 3000, end: 300}, // magenta to orange
    {start: 900, end: 2400}, // green to blue
    {start: 0, end: 0},      // all red
}
var ticksBetween = []int{4, 10, 30, 100} // number of ticks between pattern change
var levelCnt = map[mode]int{
    colorMode:     len(palettes),
    lightnessMode: len(lightnessLevels),
    speedMode:     len(ticksBetween),
}
var cur mode = colorMode
var tick time.Duration = 10 * time.Millisecond

var colors1 = make([]hsl, numPixels)
var colors2 = make([]hsl, numPixels)
var colors = make([]color.RGBA, numPixels)

func main() {
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
    palette := palettes[0]
    lightness := lightnessLevels[0]
    idx := 0
    inc := 1
    ticker := time.NewTicker(tick)
    maxTicks := ticksBetween[0]
    ticks := 0 // progress through the current color transition, minor ticks
    for {
        select {
        case config := <-configCh:
            switch config.mode {
            case colorMode:
                palette = palettes[config.level]
            case speedMode:
                maxTicks = ticksBetween[config.level]
            case lightnessMode:
                lightness = lightnessLevels[config.level]
            }
        case <-ticker.C:
            ticks = nextMinorTick(ticks, maxTicks)
            if ticks == 0 {
                // next major tick, next colors
                idx, inc = nextIndex(idx, inc)
                copy(colors1, colors2)
                updateColors(colors2, idx, inc, palette, lightness)
            }
            writeColors(colors1, colors2, int(ticks), int(maxTicks))
        }
    }
}

func nextMinorTick(ticks, maxTicks int) int {
    if ticks >= maxTicks {
        return 0
    }
    return ticks + 1
}

func writeColors(colors1, colors2 []hsl, ticks, maxTicks int) {
    for i := range numPixels {
        c1 := colors1[i]
        c2 := colors2[i]
        h := c1.H + (c2.H-c1.H)*ticks/maxTicks
        l := c1.L + (c2.L-c1.L)*ticks/maxTicks
        s := c1.S + (c2.S-c1.S)*ticks/maxTicks
        colors[i] = hsl2rgb(h, s, l)
    }
    colorWriter.WriteColors(colors)
}

func updateColors(colors2 []hsl, idx int, inc int, palette palette, l []int) {
    for i := range numPixels {
        dist := i - idx
        if dist < 0 {
            dist = -dist
        }
        l := getLightness(dist, l)
        colors2[i].H = getHue(palette.start, palette.end, int(i))
        colors2[i].S = 1000
        colors2[i].L = colorWriter.adjustLightness(l)

    }
}

func nextIndex(idx int, inc int) (int, int) {
    idx += inc
    if idx == 0 || idx == numPixels-1 {
        inc = -inc
    }
    return idx, inc
}

func colorsOff() error {
    for i := range colors {
        colors[i] = color.RGBA{}
    }
    return colorWriter.WriteColors(colors)

}

func getHue(start, end, idx int) int {
    if start > end {
        end = end + 3600
    }
    return (start + (end-start)*idx/numPixels) % 3600
}

func getLightness(dist int, lightness []int) int {
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
func hsl2rgb(h, s, l int) color.RGBA {
    a := s * min(l, 1000-l) / 1000
    f := func(n int) uint8 {
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
