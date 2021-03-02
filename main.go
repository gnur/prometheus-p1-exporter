package main

import (
	"bufio"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarm/serial"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	reader    *bufio.Reader
	powerDraw = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "power_draw_watts",
		Help: "Current power draw in Watts",
	})

	powerTariff1 = prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "power_meter_tariff1_watthours",
			Help: "power meter tariff1 reading in Watthours",
		},
		func() float64 {
			fmt.Println("reading powerTariff1Meter")
			return powerTariff1Meter
		},
	)
	powerTariff2 = prometheus.NewCounterFunc(
		prometheus.CounterOpts{
			Name: "power_meter_tariff2_watthours",
			Help: "power meter tariff2 reading in Watthours",
		},
		func() float64 {
			return powerTariff2Meter
		},
	)

	powerTariff1Meter float64
	powerTariff2Meter float64
)

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9602").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		serialPort = kingpin.Flag(
			"serial.port",
			"Serial port for the connection to the P1 interface.",
		).Required().String()
	)

	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	if *serialPort != "" {
		fmt.Println("gonna use serial device")
		config := &serial.Config{Name: *serialPort, Baud: 115200}

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

	registry := prometheus.NewRegistry()

	registry.MustRegister(powerDraw)
	registry.MustRegister(powerTariff1)
	registry.MustRegister(powerTariff2)

	fmt.Println("now serving metrics")
	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
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
			powerTariff1Meter = tmpVal * 1000
		} else if strings.HasPrefix(line, "1-0:1.8.2") {
			tmpVal, err := strconv.ParseFloat(line[10:20], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerTariff2Meter = tmpVal * 1000
		} else if strings.HasPrefix(line, "1-0:1.7.0") {
			tmpVal, err := strconv.ParseFloat(line[10:16], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			powerDraw.Set(tmpVal * 1000)
		}
	}
}
