# Pico 💓 Neopixels

Control a [Neopixel] LED strip (ws2811) with a [Raspberry Pi Pico].

https://github.com/user-attachments/assets/2db3dbb8-3ceb-4bae-b0a8-18d5291b8cb8

[Neopixel]: https://www.adafruit.com/product/2841
[Raspberry Pi Pico]: https://www.raspberrypi.com/products/raspberry-pi-pico/

## 🌈 LED emulator

Quickly prototype your LED strip patterns with a terminal-based emulator.

<details>
	<summary>🎥 Screen recording</summary>
	<img width="400" alt="screen cast terminal" src="img/terminal-led.gif">
</details>

Alternatively, try prototyping in your browser with [Evy], which can easily be
translated to Python or Go code.

[Evy]: https://play.evy.dev/#content=H4sIAAAAAAAAA31RQW7CMBC8+xUjTi0UcAC3EhKnwo0fIA7GWYhV41R2Unh+tU4cpZcmh3g9M7szG+NIB0wuTpuvifDt/XjYR2x3WEnh7K1qPMVUn5TEWmElUUgoyP/eszDkGwqsk1gukTpRbKB9iUvI1fGwf4P16NjCesOKQohHZR2hCS0JALjWAY5KBoP2N0LvM4H8RGo+a9ez+m4ZK4N+HKlkKN2R7769xV0+zGC9SYi9DiD7r8OozBuaoxgmsO8d5lmeB0RH9A25kErwlbi23vxxuvXtvW/NxySqOCSnWA6jpli/ywSW5BrNBH2JTJqPorokzH/slKjnLmia91JFhwqFlHCvI0Oj9QwmntxLYYbN4kNhOqzuXv8QnlCdG2ODcYRi0eX7BVO+g1RKAgAA

#### [Rainbow Strip]

<details>
	<summary>🎥 Screen recording</summary>
<img width="200" alt="screen cast of rainbow strip animation" src="img/rainbow-strip.gif">
</details>

[Rainbow Strip]: https://play.evy.dev/#content=H4sIAAAAAAAAA31RQW7CMBC8+xUjTi0UcAC3EhKnwo0fIA7GWYhV41R2Unh+tU4cpZcmh3g9M7szG+NIB0wuTpuvifDt/XjYR2x3WEnh7K1qPMVUn5TEWmElUUgoyP/eszDkGwqsk1gukTpRbKB9iUvI1fGwf4P16NjCesOKQohHZR2hCS0JALjWAY5KBoP2N0LvM4H8RGo+a9ez+m4ZK4N+HKlkKN2R7769xV0+zGC9SYi9DiD7r8OozBuaoxgmsO8d5lmeB0RH9A25kErwlbi23vxxuvXtvW/NxySqOCSnWA6jpli/ywSW5BrNBH2JTJqPorokzH/slKjnLmia91JFhwqFlHCvI0Oj9QwmntxLYYbN4kNhOqzuXv8QnlCdG2ODcYRi0eX7BVO+g1RKAgAA

#### [Blue-To-Magenta Ring]

<details>
	<summary>🎥 Screen recording</summary>
<img width="200" alt="screen cast of ring animaion" src="img/blue-magenta-ring.gif">
</details>

[Blue-To-Magenta Ring]: https://play.evy.dev/#content=H4sIAAAAAAAAA1WRzW6rMBCF9zzFUaVK+VGDPbapWymr2+7yBlEWBKYFXYdENtwmb39lQ1DKxp7jM2c+48px6fF0dGX19ynrhtPu8yPgfQsSmWu/m77jkOq9EdAW2kAraAFloQyUghIgCzIgBRKQFtJAKkgBCwN1yJqBUwJpAdIFyBDIWFChQa8C9FqALIGsBb1pKCHSmuqox/Poi/7YF/t1ccgclzX7mC2Q50i8HHqUXY2jv1e7z4/s0kaX2kgtzRsVJvtpWsfo/cAZAHydPRzX0eTL7psx/Yl0GL/A/Z+zm1zj3Pms9uXPLsl10rgb1wlvi8W0W0Mu8fwrOzjmC8RGZrEr+xq66tes9244TUFxm3qaiNkMvHdcH5JSs+vLqJbHME97Sazr+7jl8/118xy9L1s3Qsa2+aX3KWkMrRLDogkODaQQcMsHyIdLz2B9wyPG4tJiBcIqIeQzQnJdo6M6h9GepFuUQts9SFdsccUKLySwhhGTb4tbTH7QTud/jCtuI3PrK8eQG5NQ/wMNcHPN3gIAAA==

#### [Red strip with Trail]

<details>
	<summary>🎥 Screen recording</summary>
<img width="200" alt="screen cast of red strip animaion" src="img/red-strip.gif">
</details>

[Red strip with Trail]: https://play.evy.dev/#content=H4sIAAAAAAAAA1VRy27CMBC8+ytGnKA04ABWJUR6Kdz4A8QhOBuwakzlJIXPr7yJXfDF3sfMzqy1pdJjdLKl/h4J1133u22DdYGFFNacL62jhuODklgprCSWEguJXCGXUEdhqazIhxaJ+RwMoqZF6SqcfIz2u+07jIMm15IXxumAyIW4X4wltL4jAQD1zcNSFYq+dGfCIImL4TTUft3s0DXMNk6neuXL+55LFefI9ffQWsTHNKFMnYrBA1PHMC4kQ54mBO0FsgiPAxpL9AM5k0sRUqLunH5Ru3bddaDmp3E63AyuyLZlMD2MzpJ+U/PADeS/RW4ueAPZgHj1yr8R4X37Z/AW/qQPNxhbckhfPEnsNjDH9IG7jy/smg2NL42FRC4l7OTJ8dP+k7tH0KMwxWr2ofCWvF1vv4QHVC9VG68tIZ8ppvsDDQ1KPJoCAAA=

## 💻 Development

The LED controller is written in [Go] and built using the Go and
[TinyGo] compilers. For details on working with MicroPython and Thonny see the
[docs](python/README.md).

To build the source code, [clone] this repository and [activate Hermit] in your
terminal.

<details>
  <summary>Hermit automatically installs tools.</summary>

### Hermit

The tools used in this repository, such as Make, Go and Node, are
automatically downloaded by [Hermit] when needed. Hermit ensures that
developers on Mac, Linux, and GitHub Actions CI use the same version of
the same tools. Cloning this repo is the only installation step
necessary.

There are two ways to use the tools in the Evy repository. You can
either prefix them with `bin/`, for example `bin/make all`. Or, you can
activate Hermit in your shell with

    . ./bin/activate-hermit

This will add the tools to your path, so you can use them without having
to prefix them with `bin/`.

You can auto-activate Hermit when changing into the `evy` source
directory by installing [Hermit shell hooks] with

    hermit shell-hooks

</details>

Then, build the source code and emulate it in the terminal with

    make go

If you want to flash your source code onto a Raspberry Pi Pico, connect it to
your computer via the USB port. then run the following command:

    make flash

If you see the error the error message

    ... unable to locate any volume: [RPI-RP2]

You may need to reset your Raspberry Pi Pico first. You will need to do this
every time before you re-flash:

1. Unplug you Raspberry Pi Pico from your computer
2. Press the `BOOTSEL` button on your Pico
3. Re-plug you Pico
4. Let go of the `BOOTSEL` button

Then, run the `make flash` command again.

If you get tired of un- and re-plugging, consider adding a reset button to your
Raspberry Pi Pico. You can find instructions on how do this below in the
[Hardware section](#reset-button).

[Go]: https://go.dev
[TinyGo]: https://tinygo.org
[WebAssembly]: https://webassembly.org
[Clone]: https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository
[activate Hermit]: https://cashapp.github.io/hermit/usage/get-started/?h=activating#activating-an-environment
[Hermit]: https://cashapp.github.io/hermit
[Hermit shell hooks]: https://cashapp.github.io/hermit/usage/shell/#shell-hooks

## 🔨 Hardware

### Reset Button

optional! just if you don't want to unplug and re-plug your Raspberry Pi Pico before every re-flash.

### Trouble shooting

make onboard # makes onboard led blink. use source from onboard/main.go
