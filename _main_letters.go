package main

import (
    "image/color"
    "machine"
    "math/rand"
    "time"
)

const (
    numPixels          = 44
    longPressThreshold = 500 * time.Millisecond
    tick               = 15 * time.Millisecond

    numLetters     = 4  // Mali
    numDrawIndices = 22 // M: 6, a: 3, l: 7, i: 6
)
const (
    paletteMode mode = iota
    patternMode
    ticksMode
    lightnessMode
    modeLen
)
const (
    letterPattern pattern = iota
    drawPattern
    slidePattern
    stillPattern
    randPattern
    patternLen
)

// TODO: try polling for button for better reliability
// TODO: write last "state" to memory and use on next start
// TODO: add letter patterns, letter-by-letter; draw-letter

// Globals to avoid allocations
var (
    colors1     = make([]hsl, numPixels)
    colors2     = make([]hsl, numPixels)
    colors      = make([]color.RGBA, numPixels)
    configRange configValueRange
    letters     []letter
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

// letter holds for a given LED index the letter index and draw index
type letter struct {
    letterIdx int // For "Mali", letterIdx 0 refers to all LEDs in M
    drawIdx   int // For "Mali", drawIdx 0  refers to first LED(s) to be lit "drawing" Mali
}
type index struct {
    peak      int // peak in trail pattern
    inc       int // 1 or -1 in trail pattern movement direction
    letterIdx int // current letter index 0..numLetters
    drawIdx   int // current draw point index 0..numDrawIndices
}

func (i *index) next() {
    i.peak += i.inc
    if i.peak == 0 || i.peak == numPixels-1 {
        i.inc = -i.inc
    }
    i.letterIdx = (i.letterIdx + 1) % numLetters
    i.drawIdx = (i.drawIdx + 1) % numDrawIndices
}

func main() {
    /*
       red := color.RGBA{R: 255, B: 255}
       yellow := color.RGBA{R: 255, G: 255}
       green := color.RGBA{G: 255}
       blue := color.RGBA{B: 255}
       dur := 400 * time.Millisecond
       colors = slices.Repeat([]color.RGBA{yellow}, numPixels)
       for true {
           for i := 0; i < numPixels; i++ {
               colors[i] = color.RGBA{}
           }
           colorWriter.WriteColors(colors)
           time.Sleep(dur)
           for i := 17; i < 37; i++ {
               colors[i] = red // M
           }
           colorWriter.WriteColors(colors)
           time.Sleep(dur)
           for i := 37; i < 40; i++ {
               colors[i] = green // a - bottom
           }
           for i := 12; i < 17; i++ {
               colors[i] = green // a - top
           }
           colorWriter.WriteColors(colors)
           time.Sleep(dur)
           for i := 40; i < 42; i++ {
               colors[i] = blue // l - bottom
           }
           for i := 6; i < 12; i++ {
               colors[i] = blue // l - top
           }
           colorWriter.WriteColors(colors)
           time.Sleep(dur)
           for i := 42; i < numPixels; i++ {
               colors[i] = yellow // i - bottom
           }
           for i := 0; i < 6; i++ {
               colors[i] = yellow // i - top
           }
           colorWriter.WriteColors(colors)
           time.Sleep(dur * 3)
       }
       /*
    */
    letters = newMaliLetters()
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

func newMaliLetters() []letter {
    // M: top start: 17,end: 36; bottom start: 36, end: 37
    // a: top start: 12,end: 17; bottom start: 37, end: 40
    // l: top start: 6, end: 12; bottom start: 40, end: 42
    // i: top start: 0, end: 6; bottom start: 42, end: 44

    letters := make([]letter, 44)
    // M
    letters[23] = letter{0, 0}
    letters[24] = letter{0, 1}
    letters[25] = letter{0, 1}
    letters[22] = letter{0, 1}
    letters[26] = letter{0, 2}
    letters[27] = letter{0, 2}
    letters[21] = letter{0, 2}
    letters[28] = letter{0, 3}
    letters[29] = letter{0, 3}
    letters[20] = letter{0, 3}
    letters[30] = letter{0, 4}
    letters[31] = letter{0, 4}
    letters[19] = letter{0, 4}
    letters[32] = letter{0, 5}
    letters[33] = letter{0, 5}
    letters[18] = letter{0, 5}
    letters[34] = letter{0, 6}
    letters[35] = letter{0, 6}
    letters[36] = letter{0, 6}
    letters[17] = letter{0, 6}
    // a
    letters[37] = letter{1, 7}
    letters[16] = letter{1, 7}
    letters[38] = letter{1, 8}
    letters[15] = letter{1, 8}
    letters[14] = letter{1, 8}
    letters[39] = letter{1, 9}
    letters[13] = letter{1, 9}
    letters[12] = letter{1, 9}
    // l
    letters[40] = letter{2, 10}
    letters[11] = letter{2, 10}
    letters[10] = letter{2, 11}
    letters[9] = letter{2, 12}
    letters[8] = letter{2, 13}
    letters[7] = letter{2, 14}
    letters[41] = letter{2, 15}
    // i
    letters[6] = letter{3, 16}
    letters[40] = letter{3, 16}
    letters[5] = letter{3, 17}
    letters[4] = letter{3, 18}
    letters[3] = letter{3, 19}
    letters[2] = letter{3, 20}
    letters[1] = letter{3, 21}
    letters[42] = letter{3, 21}
    letters[0] = letter{3, 22}
    letters[43] = letter{3, 22}

    return letters
}

func startLEDs(configCh chan config) {
    idx := &index{inc: 1}
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
                idx.next()
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

func updateColors(colors1, colors2 []hsl, idx *index, config config) {
    copy(colors1, colors2)
    peak := idx
    pattern := config.pattern
    if pattern == randPattern {
        peak = rand.Intn(numPixels)
    }
    for i := range numPixels {
        l := config.lightness[0]
        if pattern == letterPattern || pattern == drawPattern {
            l = getLetterLightness(i, idx, pattern, config.lightness)
        } else if pattern != stillPattern {
            //dist := abs(i - peak)
            dist := trailDist(i, idx)
            l = getLightness(dist, config.lightness)
        }
        if pattern == stillPattern || pattern == randPattern {
            colors2[i].H = getHue(config.palette.start, config.palette.end, idx)
        } else if pattern == letterPattern {
            colors2[i].H = getLetterHue(config.palette.start, config.palette.end, i)
        } else {
            colors2[i].H = getHue(config.palette.start, config.palette.end, i)
        }
        colors2[i].S = 1000
        colors2[i].L = colorWriter.adjustLightness(l)
    }
    return
}

func getLetterLightness(i int, idx *index, pattern pattern, lightness []int) int {
    var li, peak int
    if pattern == letterPattern {
        li = letters[i].letterIdx
        peak = letters[idx.peak].letterIdx
    } else {
        li = letters[i].drawIdx
        peak = letters[idx.peak].drawIdx
    }
    dist := trailDist(peak, li, 1)
    return getLightness(dist, lightness)
}

func getLetterHue(start, end, i int, pattern pattern) int {
    if start > end {
        end += 3600
    }
    if pattern == letterPattern {
        letterIdx := letters[i].letterIdx
        return (start + (end-start)*letterIdx/numLetters) % 3600
    }
    drawIdx := letters[i].drawIdx
    return (start + (end-start)*drawIdx/numDrawIndices) % 3600
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
