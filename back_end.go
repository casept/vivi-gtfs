package main

import (
	"encoding/json"
	"log"
	"time"
)

// Generated via https://mholt.github.io/json-to-go/, matches API as of 27.04.2024
type BackEnd []struct {
	TrainColor struct {
		Color      string `json:"color"`
		Svg        string `json:"svg"`
		LoadingSvg string `json:"loadingSvg"`
		FuelType   string `json:"fuelType"`
	} `json:"trainColor"`
	StopCoordArray [][]string `json:"stopCoordArray"`
	ReturnValue    struct {
		PixelsTraveled int    `json:"pixelsTraveled"`
		I              int    `json:"i"`
		ID             string `json:"id"`
		Train          string `json:"train"`
		StopObjArray   []struct {
			Coords      []float64 `json:"coords"`
			WorkingTime []struct {
				Type    string   `json:"type"`
				Periods []string `json:"periods"`
			} `json:"workingTime"`
			ID           string `json:"_id"`
			Title        string `json:"title"`
			Departure    string `json:"departure"`
			ID0          string `json:"id"`
			GpsID        string `json:"gps_id"`
			RoutesID     string `json:"routes_id"`
			WaitingRoom  string `json:"waitingRoom"`
			Wc           string `json:"wc"`
			CoffeMachine string `json:"coffeMachine"`
			StationNotes any    `json:"stationNotes"`
			Adress       string `json:"adress"`
			I            int    `json:"i"`
		} `json:"stopObjArray"`
		AnimatedCoord []float64 `json:"animatedCoord"`
		NextTime      int       `json:"nextTime"`
		Stopped       bool      `json:"stopped"`
		CurrentStop   []float64 `json:"currentStop"`
		Position      []float64 `json:"position"`
		WaitingTime   int       `json:"waitingTime"`
		ArrivingTime  int       `json:"arrivingTime"`
		NextStopObj   struct {
			Coords      []float64 `json:"coords"`
			WorkingTime []struct {
				Type    string   `json:"type"`
				Periods []string `json:"periods"`
			} `json:"workingTime"`
			ID           string `json:"_id"`
			Title        string `json:"title"`
			Departure    string `json:"departure"`
			ID0          string `json:"id"`
			GpsID        string `json:"gps_id"`
			RoutesID     string `json:"routes_id"`
			WaitingRoom  string `json:"waitingRoom"`
			Wc           string `json:"wc"`
			CoffeMachine string `json:"coffeMachine"`
			StationNotes string `json:"stationNotes"`
			Adress       string `json:"adress"`
		} `json:"nextStopObj"`
		CurrentStopIndex int       `json:"currentStopIndex"`
		UpdaterTimeStamp time.Time `json:"updaterTimeStamp"`
		DepartureTime    string    `json:"departureTime"`
		ArrivalTime      string    `json:"arrivalTime"`
		Finished         bool      `json:"finished"`
		CurrentI         int       `json:"currentI"`
		IsGpsActive      bool      `json:"isGpsActive"`
	} `json:"returnValue"`
	Name string `json:"name"`
}

func transformBackEndMsgToGtfsRt(raw []byte) []byte {
	backEnd := BackEnd{}
	json.Unmarshal(raw, &backEnd)

	/*
	 * The data format is very strange - essentially, we constantly get a list of all stations with some train data mixed in.
	 * Unfortunately, this includes data of all trains - not just ones that are currently running!
	 */
	data := make([]GtfsRtUpdateData, 0, 10)
	for _, datum := range backEnd {
		trainNum := datum.ReturnValue.Train
		GPSActive := datum.ReturnValue.IsGpsActive
		delay := datum.ReturnValue.ArrivingTime
		route := datum.Name
		lat := datum.ReturnValue.AnimatedCoord[0]
		lon := datum.ReturnValue.AnimatedCoord[1]
		if config.verbose {
			log.Printf("Train: %s, GPS: %t, delay: %d min, route: %s, lat: %f, lon: %f\n",
			trainNum, GPSActive, delay, route, lat, lon)
		}


		tripId := lookupTripId(time.Now(), datum.ReturnValue.Train)
		if tripId == "" {
			log.Printf("Failed to find matching trip ID for back end datum %v, not including in GTFS-RT update\n", datum)
		}
		transformedDatum := GtfsRtUpdateData{
			delaySecs: int32(60 * delay),
			lat:       float32(lat),
			lon:       float32(lon),
			tripId:    tripId,
		}
		data = append(data, transformedDatum)
	}
	if len(data) > 0 {
		return createGtfsRtMsg(data)
	} else {
		log.Println("No usable data, skipping GTFS-RT update")
		return nil
	}
}
