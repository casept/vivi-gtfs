package main

import (
	"log"
	"strings"

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

func lookupTripId(serviceId string) string {
	for _, trip := range g.Trips {
		/*
		 * Service IDs usually contain the train number, but sometimes in strange forms with other junk appended.
		 * Try to work around this by counting substring matches as sufficient.
		 */
		if strings.Contains(trip.ServiceID, serviceId) {
			// FIXME: A service does not have a unique trip ID, need to look it up based on current time somehow
			return trip.ID
		}
	}
	return ""
}
func lookupGtfsStopIdFromName(name string) string {
	// Some stop names differ between GTFS feed and realtime data
	stopNameFixups := map[string]string{
		"Tukums 1":  "Tukums I",
		"Tukums 2":  "Tukums II",
		"BA Turība": "Biznesa Augstskola Turība",
		"Rēzekne 2": "Rēzekne II",
	}
	if fixedName, ok := stopNameFixups[name]; ok {
		name = fixedName
	}

	for _, stop := range g.Stops {
		if strings.ToLower(stop.Name) == strings.ToLower(name) {
			return stop.ID
		}
	}
	return ""
}
