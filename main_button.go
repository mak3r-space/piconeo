package main

import (
    "image/color"
    "machine"
    "math/rand"
    "time"
)

const isRing = true
const isLetter = true

const (
    numPixels          = 227
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
    drawPattern pattern = iota
    letterPattern
    slidePattern
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

type index struct {
    peakIdx  int
    colorIdx int // peak lightness and color in slidePattern; color in randPattern and stillPatttern
    inc      int // 1 or -1 in trail pattern movement direction

    letterIdx      int
    drawIdx        int
    numLetters     int
    numDrawIndices int
}

func (i *index) next(pattern pattern) {
    switch pattern {
    case letterPattern:
        i.letterIdx += 1
        if i.letterIdx >= i.numLetters+2 { // ad a wait add end
            i.letterIdx = -2 // add a pause of no lights
        }
    case drawPattern:
        i.drawIdx += 1
        if i.drawIdx >= i.numDrawIndices+4 { // ad a wait add end
            i.drawIdx = -2 // add a pause of no lights
        }
    case stillPattern:
        i.slideNextColor()
    case slidePattern:
        i.slideNextColor()
        i.peakIdx = i.colorIdx
    case randPattern:
        i.slideNextColor()
        i.peakIdx = rand.Intn(numPixels)
    }
}

func (i *index) slideNextColor() {
    if isRing { // 0, 1, 2, 0, 1, 2
        i.colorIdx = (i.colorIdx + numPixels - 1) % numPixels
        return
    } // 0, 1, 2, 1, 0, 1, 2
    i.colorIdx += i.inc
    if i.colorIdx == 0 || i.colorIdx == numPixels-1 {
        i.inc = -i.inc
    }
}

// letter holds for a given LED index the letter index and draw index
type letter struct {
    letterIdx int // For "Mali", letterIdx 0 refers to all LEDs in M
    drawIdx   int // For "Mali", drawIdx 0  refers to first LED(s) to be lit "drawing" Mali
}

type letterConfig struct {
    numLetters     int
    numDrawIndices int
    startIdx       int
    letters        []letter
}

func main() {
    letterConfig :=  newMeetalLetters()
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
                go startLEDs(configCh, letterConfig)
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
        paletteLevels:      newPalettes(),
        ticksBetweenLevels: newTicks(numPixels),
        lightnessLevels:    newLightnessLevels(numPixels),
    }
}

func newTicks(numPixels int) []int {
    if numPixels < 12 {
        return []int{30, 100, 200, 10} // ticks between color changes / next pattern state
    }
    return []int{4, 10, 30, 100}
}

func newLightnessLevels(numPixels int) [][]int {
    if numPixels < 12 && !isRing {
        return [][]int{
            {500, 200},
            {500},
            {20, 10},
            {20},
        }
    }
    if numPixels < 12 {
        return [][]int{
            {500, 350, 200, 100},
            {500, 200},
            {50, 20, 10, 5},
            {50, 10},
        }
    }
     if numPixels < 200 {
    return [][]int{
        {500, 350, 200, 100, 50},
        {500, 200, 50},
        {80, 50, 20, 10, 5},
        {80, 20, 5},
    }
}
    return [][]int{
        {500, 475, 450, 425, 400, 350, 200, 250, 100, 50, 20},
        {80, 75, 70, 65, 60, 55,  50, 40, 30, 20, 10, 5},
        {80, 20, 5},
    }
}

func newPalettes() []palette {
    if isLetter {
        return []palette{
            {start: 3500, end: 300},  // magenta to orange
            {start: 2000, end: 3500},  // cyan to blue
            {start: 0, end: 3600},    // rainbow
            {start: 3400, end: 3400}, // all magenta
        }
    }
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
func startLEDs(configCh chan config, letterConfig letterConfig) {
    idx := &index{
        colorIdx:       letterConfig.startIdx,
        peakIdx:        letterConfig.startIdx,
        inc:            1,
        numLetters:     letterConfig.numLetters,
        numDrawIndices: letterConfig.numDrawIndices,
    }
    ticker := time.NewTicker(tick)
    ticks := 0 // progress through the current color transition, minor ticks
    config := newStartConfig()
    pressSignal := 0
    var ok bool
    updateColors(colors1, colors2, idx, config, letterConfig)
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
                idx.next(config.pattern)
                updateColors(colors1, colors2, idx, config, letterConfig)
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

func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

func updateColors(colors1, colors2 []hsl, idx *index, config config, letterConfig letterConfig) {
    copy(colors1, colors2)
    for i := range numPixels {
        colors2[i].H = getHue(i, idx, config, letterConfig)
        colors2[i].S = 1000
        colors2[i].L = getLightness(i, idx, config.lightness, config.pattern, letterConfig)
    }
    return
}

func getHue(i int, idx *index, config config, letterConfig letterConfig) int {
    p := config.pattern
    start := config.palette.start
    end := config.palette.end
    if p == letterPattern {
        l := letterConfig.letters[i]
        return getInterpolatedHue(start, end, l.letterIdx, letterConfig.numLetters)
    }
    if p == drawPattern {
        l := letterConfig.letters[i]
        return getInterpolatedHue(start, end, l.drawIdx, letterConfig.numDrawIndices)
    }
    return getInterpolatedHue(start, end, idx.colorIdx, numPixels)
}

func getInterpolatedHue(start, end, idx, cnt int) int {
    if start > end {
        end += 3600
    }
    return (start + (end-start)*idx/cnt) % 3600
}

func getLightness(i int, idx *index, lightness []int, pattern pattern, letterConfig letterConfig) int {
    if pattern == stillPattern {
        return lightness[0]
    }
    if pattern == letterPattern {
        l := letterConfig.letters[i]
        if l.letterIdx <= idx.letterIdx {
            return lightness[0]
        } else {
            return 0
        }
    }
    if pattern == drawPattern {
        l := letterConfig.letters[i]
        if l.drawIdx <= idx.drawIdx {
            return lightness[0]
        } else {
            return 0
        }
    }
    dist := trailDist(i, idx)
    if dist >= 0 && dist < len(lightness) {
        return lightness[dist]
    }
    return 0
}

func trailDist(i int, idx *index) int {
    peak := idx.peakIdx
    if isRing {
        if peak > i {
            return peak - i
        } else {
            return peak + numPixels - i
        }
    }

    if idx.inc < 0 {
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
