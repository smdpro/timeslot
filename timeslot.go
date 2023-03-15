package timeslot

import (
	"fmt"
	"time"
)
type SlotData struct {
	Start      time.Time
	End        time.Time
	Step       int
	Period     int
	DaysInWeek []int
	Gap        int
}

var defaultDate = SlotData{
	Start:      time.Now().UTC(),
	End:        time.Now().UTC().Add(time.Hour),
	Step:       5,
	Period:     1,
	DaysInWeek: []int{1, 2, 3, 4, 5, 6, 7},
	Gap:        5,
}

func SlotGenerator(data SlotData) {

	var start, end time.Time
	if data.Start == nil {
		start = defaultDate.Start
	} else {
		start = data.Start
	}

	if data.End == nil {
		end = defaultDate.End
	} else {
		end = data.End
	}


	step := time.Duration(defaultDate.Step + defaultDate.Gap)
	current, endTime := start, start
	pm2am := false
	for current.Compare(et) < 0 {
		fmt.Println(int(current.Weekday()))

		current = time.Date(current.Year(), current.Month(), current.Day(), start.Hour(), start.Minute(), 0, 0, time.UTC)
		endTime = time.Date(current.Year(), current.Month(), current.Day(), end.Hour(), end.Minute(), 0, 0, time.UTC)
		if start.Hour() > end.Day() {
			pm2am = true
			endTime = time.Date(current.Year(), current.Month(), current.Day(), 23, 59, 59, 999, time.UTC)
		}

		if !contains(defaultDate.DaysInWeek, int(current.Weekday())) {
			current = current.AddDate(0, 0, 1)
			continue
		}
		for current.Compare(endTime) < 0 {

			fmt.Println(current)

			current = current.Add(time.Minute * step)

		}
		if pm2am {
			endTime = time.Date(current.Year(), current.Month(), current.Day(), end.Hour(), end.Minute(), 0, 0, time.UTC)
			for current.Compare(endTime) < 0 {

				fmt.Println(current)

				current = current.Add(time.Minute * step)

			}
		}
		if !pm2am {
			current = current.AddDate(0, 0, 1)
		}
	}

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
