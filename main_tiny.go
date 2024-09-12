//go:build tinygo

package main

import (
    "fmt"
    "machine"

    "tinygo.org/x/drivers/ws2812"
)

var colorWriter = newLedStrip()

type ledStrip struct {
    ws2812.Device
}

func newLedStrip() *ledStrip {
    //device := ws2812.NewWS2812(machine.GP16)
    device := ws2812.NewWS2812(machine.GP28)
    strip := &ledStrip{Device: device}
    pinConfig := machine.PinConfig{Mode: machine.PinOutput}
    strip.Pin.Configure(pinConfig)
    return strip
}

func (ls *ledStrip) adjustLightness(l int) int {
    switch {
    case l == 0:
        return 0
    case l < 10:
        return 2
    case l < 20:
        return 3
    case l < 50:
        return 4
    case l < 65:
        return 5
    case l < 75:
        return 6
    case l < 85:
        return 7
    case l < 90:
        return 8
    case l < 95:
        return 9
    case l < 100:
        return 10
    case l < 200: // 100..200 -> 10 .. 100
        return l*10/9 - 80
    case l < 350: // 200..350  -> 100 ..250
        return l - 100
    case l <= 500: // 350...500 -> 250..500
        return 250 + (l - 350)
    case l <= 1000: // mirror behavior for 500..1000
        return 1000 - ls.adjustLightness(1000-l)
    default:
        panic(fmt.Sprintf("invalid lightness %d, should be 0-1000", l))
    }
}
