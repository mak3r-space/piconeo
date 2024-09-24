package main

import (
    "image/color"
    "machine"
    "math/rand"
    "time"
)

const (
    numPixels          = 40 // 44
    longPressThreshold = 500 * time.Millisecond
    tick               = 15 * time.Millisecond
)
const (
    paletteMode mode = iota
    patternMode
    ticksMode
    lightnessMode
    modeLen
)
const (
    slidePattern pattern = iota
    stillPattern
    randPattern
    patternLen
)

// Globals to avoid allocations
var (
    colors1     = make([]hsl, numPixels)
    colors2     = make([]hsl, numPixels)
    colors      = make([]color.RGBA, numPixels)
    configRange configValueRange
)

type mode int
type pattern int

// var colorWriter ColorWriter - defined in main_tiny.go and main_biggy.go
type ColorWriter interface {
    WriteColors([]color.RGBA) error
    adjustLightness(int) int
}

type hsl struct{ H, S, L int }
type palette struct{ start, end int }
type config struct {
    palette      palette
    pattern      pattern
    ticksBetween int
    lightness    []int
    modeChange   mode
}
type configValueRange struct {
    paletteLevels      []palette
    ticksBetweenLevels []int
    lightnessLevels    [][]int
}

func main() {
    configRange = newConfigValueRange(numPixels)
    levels := make(map[mode]int, modeLen)
    levelLens := newLevelLens()
    cur := mode(0)

    levelCh := makeButtonChan(machine.GP0)
    modeCh := makeButtonChan(machine.GP15)
    var configCh chan config
    off := true
    modeCh <- true
    for {
        select {
        case longPress := <-modeCh:
            switch {
            case longPress && off: // initialize and turn on again
                cur = 0
                resetLevels(levels)
                configCh = make(chan config, 32)
                go startLEDs(configCh)
                off = false
            case longPress && !off: // turn off
                close(configCh)
                off = true
            default:
                cur = (cur + 1) % modeLen
                configCh <- newModeChangeConfig(levels, cur) // just for press signal
            }
        case <-levelCh:
            if !off {
                levels[cur] = (levels[cur] + 1) % levelLens[cur]
                configCh <- newConfig(levels)
            }
        }
    }
}

func resetLevels(levels map[mode]int) {
    for mode := range levels {
        levels[mode] = 0
    }
}

func newConfig(levels map[mode]int) config {
    return newModeChangeConfig(levels, modeLen)
}

func newModeChangeConfig(levels map[mode]int, m mode) config {
    return config{
        palette:      configRange.paletteLevels[levels[paletteMode]],
        pattern:      pattern(levels[patternMode]),
        ticksBetween: configRange.ticksBetweenLevels[levels[ticksMode]],
        lightness:    configRange.lightnessLevels[levels[lightnessMode]],
        modeChange:   m,
    }
}

func newStartConfig() config {
    return config{
        palette:      configRange.paletteLevels[0],
        pattern:      pattern(0),
        ticksBetween: configRange.ticksBetweenLevels[0],
        lightness:    configRange.lightnessLevels[0],
        modeChange:   modeLen,
    }
}

func newConfigValueRange(numPixels int) configValueRange {
    return configValueRange{
        paletteLevels:      newPallettes(),
        ticksBetweenLevels: newTicks(numPixels),
        lightnessLevels:    newLightnessLevels(numPixels),
    }
}

func newTicks(numPixels int) []int {
    if numPixels < 12 {
        return []int{30, 100, 200, 10} // ticks between color changes / next pattern state
    }
    return []int{10, 30, 100, 4}
}

func newLightnessLevels(numPixels int) [][]int {
    if numPixels < 12 {
        return [][]int{
            {500, 200},
            {500},
            {20, 10},
            {20},
        }
    }
    return [][]int{
        {500, 350, 200, 100, 50},
        {500, 200, 50},
        {80, 50, 20, 10, 5},
        {80, 20, 5},
    }
}

func newPallettes() []palette {
    return []palette{
        {start: 0, end: 3600},   // rainbow
        {start: 3000, end: 300}, // magenta to orange
        {start: 900, end: 2400}, // green to blue
        {start: 0, end: 0},      // all red
    }
}

func newLevelLens() map[mode]int {
    return map[mode]int{
        paletteMode:   len(configRange.paletteLevels),
        patternMode:   int(patternLen),
        ticksMode:     len(configRange.ticksBetweenLevels),
        lightnessMode: len(configRange.lightnessLevels),
    }

}
func startLEDs(configCh chan config) {
    idx := 0
    inc := 1
    ticker := time.NewTicker(tick)
    ticks := 0 // progress through the current color transition, minor ticks
    config := newStartConfig()
    pressSignal := 0
    var ok bool
    for {
        select {
        case config, ok = <-configCh:
            if !ok {
                ledsOff()
                return // channel closed
            }
            pressSignal = 15
        case <-ticker.C:
            if pressSignal > 0 {
                signalLEDs(config.modeChange)
                pressSignal -= 1
                continue
            }
            ticks = nextMinorTick(ticks, config.ticksBetween)
            if ticks == 0 {
                // next major tick, next colors
                idx, inc = nextIndex(idx, inc, config.pattern)
                updateColors(colors1, colors2, idx, inc, config)
            }
            writeColors(colors1, colors2, colors, ticks, config.ticksBetween)
        }
    }
}

func nextMinorTick(ticks, maxTicks int) int {
    if ticks >= maxTicks {
        return 0
    }
    return ticks + 1
}

func writeColors(colors1, colors2 []hsl, colors []color.RGBA, ticks, maxTicks int) {
    for i := range numPixels {
        c1, c2 := colors1[i], colors2[i]
        h1 := c1.H
        h2 := c2.H
        if abs(h1-h2) > 1800 {
            if h1 < h2 {
                h1 += 3600
            } else {
                h2 += 3600
            }
        }
        h := (h1 + (h2-h1)*ticks/maxTicks) % 3600
        l := c1.L + (c2.L-c1.L)*ticks/maxTicks
        s := c1.S + (c2.S-c1.S)*ticks/maxTicks
        colors[i] = hsl2rgb(h, s, l)
    }
    colorWriter.WriteColors(colors)
}

func updateColors(colors1, colors2 []hsl, idx int, inc int, config config) {
    copy(colors1, colors2)
    peak := idx
    pattern := config.pattern
    if pattern == randPattern {
        peak = rand.Intn(numPixels)
    }
    for i := range numPixels {
        l := config.lightness[0]
        if pattern != stillPattern {
            //dist := abs(i - peak)
            dist := trailDist(peak, i, inc)
            l = getLightness(dist, config.lightness)
        }
        if pattern == stillPattern || pattern == randPattern {
            colors2[i].H = getHue(config.palette.start, config.palette.end, idx)
        } else {
            colors2[i].H = getHue(config.palette.start, config.palette.end, i)
        }
        colors2[i].S = 1000
        colors2[i].L = colorWriter.adjustLightness(l)
    }
    return
}

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

func trailDist(peak, i, inc int) int {
    if inc < 0 {
        return i - peak
    }
    return peak - i
}
func nextIndex(idx int, inc int, pattern pattern) (int, int) {
    idx += inc
    if idx == 0 || idx == numPixels-1 {
        inc = -inc
    }
    return idx, inc
}

func ledsOff() error {
    for i := range colors {
        colors[i] = color.RGBA{}
    }
    return colorWriter.WriteColors(colors)
}

func signalLEDs(m mode) {
    var c color.RGBA
    switch m {
    case paletteMode:
        c = color.RGBA{R: 255}
    case patternMode:
        c = color.RGBA{R: 255, G: 255}
    case ticksMode:
        c = color.RGBA{G: 255}
    case lightnessMode:
        c = color.RGBA{B: 255}
    case modeLen: // no modeChange
        c = color.RGBA{R: 255, G: 255, B: 255}

    }
    for i := range colors {
        colors[i] = c
    }
    colorWriter.WriteColors(colors)
}

func getHue(start, end, idx int) int {
    if start > end {
        end += 3600
    }
    return (start + (end-start)*idx/numPixels) % 3600
}

func getLightness(dist int, lightness []int) int {
    if dist >= 0 && dist < len(lightness) {
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

func makeButtonChan(btn machine.Pin) chan bool {
    in := setupButtonPressChan(btn)
    out := make(chan bool, 32)
    go handleButtonPress(in, out)
    return out
}

func handleButtonPress(in <-chan event, out chan bool) {
    start := time.Now() // press start time
    lastEvent := start  // last press or release event
    for e := range in { // button down (true) or up (false) event
        now := e.t
        if now.Sub(lastEvent) < 5*time.Millisecond {
            continue
        }
        lastEvent = now
        if e.press {
            start = now
        } else {
            out <- now.Sub(start) > longPressThreshold // long/short press
        }
    }
}

type event struct {
    t     time.Time
    press bool
}

func setupButtonPressChan(btn machine.Pin) chan event {
    config := machine.PinConfig{Mode: machine.PinInputPullup}
    ch := make(chan event, 32)
    btn.Configure(config)
    btn.SetInterrupt(machine.PinFalling|machine.PinRising, func(pin machine.Pin) {
        select {
        case ch <- event{t: time.Now(), press: !pin.Get()}:
        default:
        }
    })
    // go pollButtonPress(btn, ch)
    return ch
}

func pollButtonPress(btn machine.Pin, ch chan event) chan event {
    ticker := time.NewTicker(10 * time.Millisecond)
    state := false
    for {
        select {
        case <-ticker.C:
            if state != btn.Get() {
                state = !state
                ch <- event{t: time.Now(), press: state}
            }
        }
    }
}
