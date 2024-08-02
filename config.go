package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	listenAddr string
	verbose bool
}

var config Config

func handleCliArgs() {
	config.listenAddr = *flag.String("listenAddr", "localhost:1337", "Address to expose the GTFS-RT HTTP server on")
	config.verbose = *flag.Bool("verbose", false, "Activate verbose debug logging")
	flag.Parse()
}

var Usage = func() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}
