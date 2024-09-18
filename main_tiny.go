//go:build tinygo

package main

import (
    "machine"

    "tinygo.org/x/drivers/ws2812"
)

var colorWriter = newLedStrip()

type ledStrip struct {
    ws2812.Device
}

func newLedStrip() *ledStrip {
    device := ws2812.NewWS2812(machine.GP28)
    strip := &ledStrip{Device: device}
    pinConfig := machine.PinConfig{Mode: machine.PinOutput}
    strip.Pin.Configure(pinConfig)
    return strip
}

func (*ledStrip) getLightness(distCenter int) float64 {
    switch distCenter {
    case 0:
        return 0.5
    case 1:
        return 0.25
    case 2:
        return 0.1
    case 3:
        return 0.01
    case 4:
        return 0.0025
    }
    return 0
}
