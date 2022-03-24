package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
)

var (
	reader    *bufio.Reader
	powerDraw = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "power_draw_watts",
			Help: "Current power draw in Watts",
		},
		[]string{"type"}, //type should either be delivered or received
	)

	powerMeter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "power_meter_watthours",
			Help: "power meter reading in Watthours",
		},
		[]string{"type", "tariff"}, //type is either delivered or received, tariff is either high or low
	)

	gasMeter = prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "gas_meter_cm2",
			Help: "Gas meter reading in cm2",
		},
		func() float64 {
			return gasTotalMeter
		},
	)

	gasTotalMeter float64
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(powerDraw)
	prometheus.MustRegister(powerMeter)
	prometheus.MustRegister(gasMeter)

}

func main() {
	if os.Getenv("SERIAL_DEVICE") != "" {
		fmt.Println("gonna use serial device")
		config := &serial.Config{Name: os.Getenv("SERIAL_DEVICE"), Baud: 115200}

		usb, err := serial.OpenPort(config)
		if err != nil {
			fmt.Printf("Could not open serial interface: %s", err)
			return
		}

		reader = bufio.NewReader(usb)
	} else {
		fmt.Println("gonna use some files")
		file, err := os.Open("examples/fulllist.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	go listener(reader)

	// sleeping 10 seconds to prevent uninitialized scrapes
	time.Sleep(10 * time.Second)

	fmt.Println("now serving metrics")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9222", nil))

}

func listener(source io.Reader) {
	var line string
	for {
		rawLine, err := reader.ReadBytes('\x0a')
		if err != nil {
			fmt.Println(err)
			return
		}
		line = string(rawLine[:])
		if strings.HasPrefix(line, "1-0:1.8.1") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerMeter.WithLabelValues("low", "received").Set(tmpVal * 1000)
		} else if strings.HasPrefix(line, "1-0:1.8.2") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerMeter.WithLabelValues("high", "received").Set(tmpVal * 1000)

		} else if strings.HasPrefix(line, "1-0:2.8.1") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerMeter.WithLabelValues("low", "delivered").Set(tmpVal * 1000)
		} else if strings.HasPrefix(line, "1-0:2.8.2") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerMeter.WithLabelValues("high", "delivered").Set(tmpVal * 1000)

		} else if strings.HasPrefix(line, "0-1:24.2.1") {
			tmpVal, err := strconv.ParseFloat(line[26:35], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			gasTotalMeter = tmpVal * 100 * 100 * 100 // m3 to cm3
		} else if strings.HasPrefix(line, "1-0:1.7.0") {
			tmpVal, err := strconv.ParseFloat(line[10:16], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerDraw.WithLabelValues("received").Set(tmpVal * 1000)
		} else if strings.HasPrefix(line, "1-0:2.7.0") {
			tmpVal, err := strconv.ParseFloat(line[10:16], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerDraw.WithLabelValues("delivered").Set(tmpVal * 1000)
		}
		if os.Getenv("SERIAL_DEVICE") == "" {
			time.Sleep(200 * time.Millisecond)
		}
	}
}
