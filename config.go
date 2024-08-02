package main

import (
	"flag"
	"os"
)

type Config struct {
	gtfsPath   string
	listenAddr string
	verbose    bool
}

var config Config

func handleCliArgs() {
	flag.StringVar(&config.gtfsPath, "gtfs-path", "", "Path to ZIP containing static GTFS data to match against")
	flag.StringVar(&config.listenAddr, "listen-addr", "localhost:1337", "Address to expose the GTFS-RT HTTP server on")
	flag.BoolVar(&config.verbose, "verbose", false, "Activate verbose debug logging")
	flag.Parse()

	if config.gtfsPath == "" {
		flag.Usage()
		os.Exit(1)
	}
}
