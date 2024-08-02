package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/artonge/go-gtfs"
)

// The live websocket data does not provide a trip ID, which is needed to generate the real-time feed.
// Here, we parse just enough static GTFS to enable mapping the current time and service ID from the feed to a trip ID.

var g *gtfs.GTFS = nil

// See https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file
// Modfified to unzip to a temporary directory and return path to it
func unzip(src string) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	dir, err := os.MkdirTemp("", "vivi-gtfs")
	if err != nil {
		return "", err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dir, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return "", err
		}
	}

	return dir, nil
}

func initGtfs(dataPath string) {
	// Extract GTFS zip file, library only supports loading from directory
	dir, err := unzip(dataPath)
	if err != nil {
		log.Fatalf("Failed to extract GTFS data for mapping trip IDs: %s\n", err)
	}
	defer os.RemoveAll(dir)

	gt, err := gtfs.Load(dir, nil)
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
