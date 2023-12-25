package main

import (
	"image/color"
	"machine"
	"strconv"
	"time"

	font "github.com/Nondzu/ssd1306_font"
	"tinygo.org/x/drivers/bme280"
	"tinygo.org/x/drivers/ssd1306"
)

const (
	LED_PIN         = machine.GPIO15
	DISPLAY_SDA_PIN = machine.GPIO6
	DISPLAY_SCL_PIN = machine.GPIO7
	BME280_SDA_PIN  = machine.GPIO0
	BME280_SCL_PIN  = machine.GPIO1
)

var (
	sensor     bme280.Device
	displayDev ssd1306.Device
	display    font.Display
)

func drawGlyph(x int16, y int16, glyph [GLYPH_HEIGHT][GLYPH_WIDTH]int16) {
	for i := int16(0); i < GLYPH_HEIGHT; i++ {
		for j := int16(0); j < GLYPH_WIDTH; j++ {
			if glyph[i][j] == 1 {
				displayDev.SetPixel(x+j, y+i, color.RGBA{255, 255, 255, 255})
			}
		}
	}
}

func writeToDisplay(x int16, y int16, text string) {
	display.YPos = y
	display.XPos = x
	display.PrintText(text)
}

func setupDisplay() {
	// Display
	machine.I2C1.Configure(machine.I2CConfig{
		Frequency: machine.KHz * 400,
		SDA:       DISPLAY_SDA_PIN,
		SCL:       DISPLAY_SCL_PIN,
	})

	displayDev = ssd1306.NewI2C(machine.I2C1)
	displayDev.Configure(ssd1306.Config{
		Address: 0x3C,
		Width:   128,
		Height:  64,
	})

	// Reset Display
	displayDev.ClearBuffer()
	displayDev.ClearDisplay()

	display = font.NewDisplay(displayDev)
	display.Configure(font.Config{FontType: font.FONT_7x10})

	// Title
	writeToDisplay(0, 0, "Weather Station")
}

func setupSensor() {
	// BME280
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: machine.KHz * 400,
		SDA:       BME280_SDA_PIN,
		SCL:       BME280_SCL_PIN,
	})

	sensor = bme280.New(machine.I2C0)
	sensor.Configure()

	connected := sensor.Connected()
	if !connected {
		println("BME280 not detected")
	}
	println("BME280 detected")

}

func main() {
	// Power Light
	LED_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	LED_PIN.High()

	setupDisplay()
	setupSensor()

	for {
		temp, _ := sensor.ReadTemperature()
		tempString := strconv.FormatFloat(((float64(temp)*1.8)/1000)+float64(32.0000), 'f', 1, 64) + " F"
		drawGlyph(0, 16, getTemperatureGlyph())
		writeToDisplay(14, 18, tempString)

		hum, _ := sensor.ReadHumidity()
		humString := strconv.FormatFloat(float64(hum)/100, 'f', 1, 64) + "%"
		drawGlyph(0, 32, getHumidityGlyph())
		writeToDisplay(14, 34, humString)

		press, _ := sensor.ReadPressure()
		pressString := strconv.FormatFloat(float64(press)/100000, 'f', 1, 64) + " hPa"
		drawGlyph(0, 48, getPressureGlyph())
		writeToDisplay(14, 51, pressString)

		alt, _ := sensor.ReadAltitude()
		println("Altitude:", alt, "m")

		time.Sleep(10 * time.Second)
	}
}
