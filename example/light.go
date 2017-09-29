package main

import (
	gpio "github.com/chenqinghe/raspberrypi-gpio"
	"time"
)

func main() {
	light := gpio.NewPin("light", 4, gpio.OUT)
	if err := light.Export(); err != nil {
		panic(err)
	}
	defer light.Unexport()
	if err := light.Write(gpio.HIGH); err != nil {
		panic(err)
	}
	if err := light.Blink(time.Second, time.Second*10); err != nil {
		panic(err)
	}
}
