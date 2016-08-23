package main

import (
	"os"
	"os/signal"
	"gate"
	"flag"
)

func main() {
	var broker string
	var udpport int

	flag.StringVar(&broker, "broker", "tcp://localhost:1883", "Configuration File")
	flag.IntVar(&udpport, "port", 1884, "MQTT-SN Gateway UDP Listening Port")
	flag.Parse()

	var gw gateway.Gateway
	stopsig := registerSignals()

	gateway.InitLogger(os.Stdout, os.Stderr)

	gw = initAggregating(udpport, broker, stopsig)

	gw.Start()
}

func initAggregating(udpport int, broker string, stopsig chan os.Signal) *gateway.AGateway {
	a := gateway.NewAGateway(udpport, broker, stopsig)
	return a
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}