from neopixel import Neopixel
import utime
import random

led_count = 20
strip = Neopixel(led_count, 0, 28, "GRB")
strip.brightness(50)
while True:
    for i in range(0, led_count):
        color = (1, 1, 1) # red, green, blue 0..255
        strip.set_pixel(i, color)
        strip.show()
        utime.sleep(0.2)
    strip.clear()
    strip.show()
    utime.sleep(1)
