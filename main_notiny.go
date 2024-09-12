//go:build !tinygo

package main

import (
    "fmt"
    "image/color"

    tty "github.com/gookit/color"
)

type terminal struct{ sep string }

var colorWriter = terminal{sep: " "}

func (t terminal) WriteColors(colors []color.RGBA) error {
    if len(colors) == 0 {
        return nil
    }
    c := colors[0]
    fmt.Print("\r")
    tty.RGB(c.R, c.G, c.B).Print("⬤")
    for _, c := range colors[1:] {
        fmt.Print(t.sep)
        tty.RGB(c.R, c.G, c.B).Print("⬤")
    }
    return nil
}

func (terminal) adjustLightness(l int32) int32 {
    return l
}
