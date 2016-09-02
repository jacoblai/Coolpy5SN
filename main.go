package main

import (
	"os"
	"os/signal"
	"gate"
	"flag"
)

func main() {
	var (
		broker = flag.String("b", "tcp://localhost:1883", "Coolpy5 mqtt broker addr default tcp://localhost:1883")
		cphttp = flag.String("h", "http://localhost:6543", "Coolpy5 http api addr default http://localhost:6543")
		udpport = flag.Int("p", 1884, "MQTT-SN Gateway UDP Listening Port default 1884")
	)
	flag.Parse()

	var gw gateway.Gateway
	stopsig := registerSignals()

	gateway.InitLogger(os.Stdout, os.Stderr)

	gw = initAggregating(*udpport, *broker, *cphttp, stopsig)

	gw.Start()
}

func initAggregating(udpport int, broker string, cphttp string, stopsig chan os.Signal) *gateway.AGateway {
	a := gateway.NewAGateway(udpport, broker, cphttp, stopsig)
	return a
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}