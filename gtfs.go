package main

import (
	"log"
	"time"

	"github.com/artonge/go-gtfs"
)

// The live websocket data does not provide a trip ID, which is needed to generate the real-time feed.
// Here, we parse just enough static GTFS to enable mapping the current time and service ID from the feed to a trip ID.

var g *gtfs.GTFS = nil

func initGtfs(dataPath string) {
	gt, err := gtfs.Load(dataPath, nil)
	if err != nil {
		log.Fatalf("Failed to load GTFS data for mapping trip IDs: %s\n", err)
	}
	g = gt
}

func lookupTripId(t time.Time, serviceId string) string {
	for _, trip := range g.Trips {
		if trip.ServiceID == serviceId {
			// FIXME: A service does not have a unique trip ID, need to look it up based on current time somehow
			return trip.ID
		}
	}
	return ""
}
