package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const URL = "wss://trainmap.pv.lv/ws"

func main() {
	handleCliArgs()

	log.Printf("Loading static GTFS data from %s\n", config.gtfsPath)
	initGtfs(config.gtfsPath)

	log.Println("Connecting to API...")
	ws, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	log.Println("Connected to API websocket")

	gtfsRtData := make(chan []byte)
	go listenAndServeFeed(gtfsRtData)

	for {
		var m map[string]interface{}
		if err = ws.ReadJSON(&m); err != nil {
			log.Fatal(err)
		}
		handleMessage(m, gtfsRtData)
	}
}

func handleMessage(m map[string]interface{}, gtfsRtData chan []byte) {
	switch fmt.Sprintf("%s", m["type"]) {
	case "message":
		{
			log.Println("Received message of type \"message\", ignoring as it's always empty")
		}
	case "active-stops":
		{
			log.Println("Received message of type \"active-stops\", probably not relevant")
		}
	case "back-end":
		{
			b, err := json.Marshal(m["data"])
			if err != nil {
				log.Fatalf("Failed to marshal back-end data: %s\n", err)
			}
			data := transformBackEndMsgToGtfsRt(b)
			if data != nil {
				gtfsRtData <- data
			}
		}
	default:
		{
			log.Println("Received message of unknown type, dumping:")
			log.Println(m)
		}
	}
}
