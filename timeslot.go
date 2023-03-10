package timeslot

import "time"

type periods struct {
	minute string
	hour   string
	day    string
}
type SlotData struct {
	Start      time.Time
	End        time.Time
	Step       int
	Period     uint
	DaysInWeek [7]uint
	Gap        uint
}

var defaultDate = SlotData{
	Start:      time.Now().UTC(),
	End:        time.Now().UTC().Add(time.Hour),
	Step:       5,
	Period:     5,
	DaysInWeek: [7]uint{1, 2, 3, 4, 5, 6, 7},
	Gap:        0,
}

func SlotGenerator(data SlotData) {

}
