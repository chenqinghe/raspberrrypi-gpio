package gpio

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	SYSFS_GPIO_EXPORT        = "/sys/class/gpio/export"
	SYSFS_GPIO_UNEXPORT      = "/sys/class/gpio/unexport"
	SYSFS_GPIO_PIN_VALUE     = "/sys/class/gpio/gpio%d/value"
	SYSFS_GPIO_PIN_DIRECTION = "/sys/class/gpio/gpio%d/direction"
)

const (
	HIGH int = 1
	LOW  int = 0
)

type mode string

const (
	IN  mode = "in"
	OUT mode = "out"
)

type Pin struct {
	Name     string //Name字段助记功能，方便知道是干什么用的
	Number   int    //此number为BCM编号的引脚值
	Mode     mode
	Mux      *sync.Mutex
	exported bool
}

func NewPin(name string, num int, mode mode) *Pin {
	mux := &sync.Mutex{}
	return &Pin{
		Name:   name,
		Number: num,
		Mode:   mode,
		Mux:    mux,
	}
}

func (pin *Pin) Export() (err error) {
	f, err := os.OpenFile(SYSFS_GPIO_EXPORT, os.O_WRONLY, 755)
	if err != nil {
		return
	}
	defer f.Close()
	//export pin
	_, err = f.Write([]byte(strconv.Itoa(pin.Number)))
	if err != nil {
		return
	}
	//set mode
	directionFile := fmt.Sprintf(SYSFS_GPIO_PIN_DIRECTION, pin.Number)
	f2, err := os.OpenFile(directionFile, os.O_WRONLY, 755)
	if err != nil {
		return
	}
	defer f2.Close()
	_, err = f2.Write([]byte(pin.Mode))
	if err != nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = errors.New("Pin.mux cannot be nil")
		}
	}()
	pin.Mux.Lock()
	pin.exported = true
	pin.Mux.Unlock()
	return
}

func (pin *Pin) Unexport() (err error) {
	f, err := os.OpenFile(SYSFS_GPIO_UNEXPORT, os.O_WRONLY, 755)
	if err != nil {
		return
	}
	defer f.Close()
	//unexport gpio
	_, err = f.Write([]byte(strconv.Itoa(pin.Number)))
	if err != nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			err = errors.New("Pin.mux cannot be nil")
		}
	}()
	pin.Mux.Lock()
	pin.exported = false
	pin.Mux.Unlock()
	return
}

func (pin *Pin) Write(lv int) error {
	if !pin.exported {
		return errors.New("cannot operate unexported pin.")
	}
	valueFile := fmt.Sprintf(SYSFS_GPIO_PIN_VALUE, pin.Number)
	f, err := os.OpenFile(valueFile, os.O_WRONLY, 755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(strconv.Itoa(lv)))
	return err
}

func (pin *Pin) Read() (int, error) {
	if !pin.exported {
		return -1, errors.New("cannot operate unexported pin.")
	}
	valueFile := fmt.Sprintf(SYSFS_GPIO_PIN_VALUE, pin.Number)
	f, err := os.OpenFile(valueFile, os.O_RDONLY, 755)
	if err != nil {
		return -1, err
	}
	defer f.Close()
	status := make([]byte, 1)
	_, err = f.Read(status)
	if err != nil {
		return -1, err
	}
	i, err := strconv.Atoi(string(status))
	if err != nil {
		return -1, err
	}
	return i, nil
}

func (pin *Pin) Blink(interval time.Duration, last time.Duration) error {
	var timeup <-chan time.Time
	if last <= 0 {
		timeup = time.After(math.MaxInt64)
	} else {
		timeup = time.After(last)
	}
	status, err := pin.Read()
	if err != nil {
		return err
	}
LOOP:
	for {
		select {
		case <-time.After(interval):
			status ^= 1
			err := pin.Write(status)
			if err != nil {
				return err
			}
		case <-timeup:
			break LOOP
		}
	}
	return nil
}

func (pin *Pin) Toggle() error {
	status, err := pin.Read()
	if err != nil {
		return err
	}
	return pin.Write(status ^ 1)
}
