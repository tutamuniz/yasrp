package main

import (
	"flag"
	"log"

	"github.com/tutamuniz/yasrp/minihttp/reverseproxy"
	"github.com/tutamuniz/yasrp/miniutils/config"
)

func main() {
	configFile := flag.String("config", "config.json", "Configuration File.(JSON format)")
	flag.Parse()

	config, err := config.ParseConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err.Error())
	}

	rp, err := reverseproxy.NewReverseProxyFromConfig(*config)
	if err != nil {
		log.Fatalln(err)
	}

	rp.Listen()
}
