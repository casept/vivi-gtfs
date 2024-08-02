package main

import (
	"log"
	"net/http"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

type GtfsRtUpdateData struct {
	delaySecs int32
	lat       float32
	lon       float32
	tripId    string
}

func addToFeed(msg *gtfs.FeedMessage, data GtfsRtUpdateData) {
	pos := gtfs.Position{
		Latitude:  &data.lat,
		Longitude: &data.lon,
	}
	descr := gtfs.TripDescriptor{
		TripId: &data.tripId,
	}
	update := gtfs.TripUpdate{
		/* TODO: TripDescriptor */
		Trip:  &descr,
		Delay: &data.delaySecs,
	}
	vPos := gtfs.VehiclePosition{
		Position: &pos,
	}

	/* Each FeedEntity should contain only one type of update, so we need 2 */

	e := gtfs.FeedEntity{
		Id:         &data.tripId,
		TripUpdate: &update,
	}

	msg.Entity = append(msg.Entity, &e)

	e2 := gtfs.FeedEntity{
		Id:      &data.tripId,
		Vehicle: &vPos,
	}
	msg.Entity = append(msg.Entity, &e2)
}

func createGtfsRtMsg(data []GtfsRtUpdateData) []byte {
	var now uint64 = uint64(time.Now().Unix())
	var gtfs_version string = "2.0"
	feed := gtfs.FeedMessage{}
	feed.Header = &gtfs.FeedHeader{
		Timestamp:           &now,
		GtfsRealtimeVersion: &gtfs_version,
	}

	for _, datum := range data {
		addToFeed(&feed, datum)
	}

	buf, err := proto.Marshal(&feed)
	if err != nil {
		log.Fatalf("Failed to serialize RTFS-RT protobuf message: %s\n", err)
	}
	if config.verbose {
		log.Printf("GTFS-RT message:\n%s\n", feed.String())
	}
	return buf
}

// Serve feed data via HTTP, updating it when new one comes in from the channel
func listenAndServeFeed(dataChan chan []byte) {
	// Block here until we have some data
	feedData := <-dataChan

	go func() {
		for {
			// FIXME: Prevent data race by locking this when updating
			feedData = <-dataChan
		}
	}()

	http.DefaultServeMux.HandleFunc("/gtfs-rt.pb", func(w http.ResponseWriter, r *http.Request) {
		w.Write(feedData)
	})
	http.ListenAndServe(config.listenAddr, nil)
}
