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
	DaysInWeek []int
	Gap        uint
}

var defaultDate = SlotData{
	Start:      time.Now().UTC(),
	End:        time.Now().UTC().Add(time.Hour),
	Step:       5,
	Period:     5,
	DaysInWeek: []int{1, 2, 3, 4, 5, 6, 7},
	Gap:        0,
}

func SlotGenerator(data SlotData) {

}


func index[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

func contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}
