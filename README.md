# raspberrrypi-gpio   [中文说明](https://github.com/chenqinghe/raspberrrypi-gpio/blob/master/READM_ZH.md)

this project is a library to operate raspberry pi gpio pins.


# usage

### example 
here is a example to control led.

```GO
ledpin := gpio.NewPin("led control", 4, gpio.OUT)

//expose pin to user space
if err := ledpin.Export(); err != nil {
	panic(err)
}
defer ledpin.Unexport()

//light led
if err := ledpin.Write(gpio.HIGH); err != nil {
	panic(err)
}

//twinkle
if err := ledpin.Blink(time.Second, time.Second*10); err != nil {
	panic(err)
}
```

# license
this project is under [MIT](https://github.com/chenqinghe/raspberrrypi-gpio/blob/master/LICENSE) license.
