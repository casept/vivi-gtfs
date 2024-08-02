package main

type ActiveStops []struct {
	Coords           []string  `json:"coords"`
	WorkingTime      []any     `json:"workingTime"`
	DirectionList    []int     `json:"directionList"`
	ID               string    `json:"_id"`
	ID0              string    `json:"id"`
	PvID             string    `json:"pvID,omitempty"`
	Title            string    `json:"title"`
	LockedTitle      string    `json:"lockedTitle,omitempty"`
	WaitingRoom      string    `json:"waitingRoom"`
	Canopy           string    `json:"canopy,omitempty"`
	Wc               string    `json:"wc"`
	CoffeMachine     string    `json:"coffeMachine"`
	StationNotes     any       `json:"stationNotes"`
	Adress           string    `json:"adress"`
	Departure        string    `json:"departure,omitempty"`
	GpsID            string    `json:"gps_id,omitempty"`
	RoutesID         string    `json:"routes_id,omitempty"`
	I                int       `json:"i,omitempty"`
	Train            string    `json:"train,omitempty"`
	CurrentStopIndex int       `json:"currentStopIndex,omitempty"`
	AnimatedCoord    []float64 `json:"animatedCoord,omitempty"`
	StopIndex        int       `json:"stopIndex,omitempty"`
}
