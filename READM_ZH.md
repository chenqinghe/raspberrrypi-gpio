# raspberrypi-gpio   [English](https://github.com/chenqinghe/raspberrrypi-gpio/blob/master/README.md)

这是一个用go实现的控制树莓派引脚输入输出的库。虽然目前github上已经有现成的库了，但是要么是通过cgo调用C语言的wiringpi库来实现的，要么实现的接口不是很友好，因此自己撸了一个出来，希望有兴趣的伙伴可以多提issue和pr.


# 使用方法

### 示例 
以下是一个控制led灯的示例

```GO
ledpin := gpio.NewPin("led control", 4, gpio.OUT)

//将此引脚暴露至用户空间
if err := ledpin.Export(); err != nil {
	panic(err)
}
defer ledpin.Unexport()

//写入高电平
if err := ledpin.Write(gpio.HIGH); err != nil {
	panic(err)
}

//让led闪烁
if err := ledpin.Blink(time.Second, time.Second*10); err != nil {
	panic(err)
}
```

# 协议
本项目遵从MIT协议，具体内容见[这里](https://github.com/chenqinghe/raspberrrypi-gpio/blob/master/LICENSE)
