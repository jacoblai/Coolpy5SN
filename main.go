package main

import (
	"os"
	"os/signal"
	G "gateway/gate"
	"flag"
)

func main() {
	var broker string
	var udpport int

	flag.StringVar(&broker, "broker", "tcp://localhost:1883", "Configuration File")
	flag.IntVar(&udpport, "port", 1884, "MQTT-SN Gateway UDP Listening Port")
	flag.Parse()

	var gateway G.Gateway
	stopsig := registerSignals()

	G.InitLogger(os.Stdout, os.Stderr) // todo: configurable

	gateway = initTransparent(udpport,broker, stopsig)

	gateway.Start()
}

func initTransparent(udpport int,broker string, stopsig chan os.Signal) *G.TGateway {
	t := G.NewTGateway(udpport,broker, stopsig)
	return t
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}