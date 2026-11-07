package timeslot

import (
	"fmt"
	"strings"
	"time"
)

// Period represents the time period unit
type Period string

const (
	PeriodMinute Period = "m"
	PeriodHour   Period = "h"
	PeriodDay    Period = "d"
)

// Options holds the configuration for slot generation
type Options struct {
	Start      interface{} // time.Time, string, or int64 (timestamp)
	End        interface{}
	Step       int
	Period     Period
	DaysInWeek []int // ISO weekdays 1-7 (Mon=1, Sun=7)
	Gap        int
}

// DefaultOptions returns the default options
func DefaultOptions() Options {
	return Options{
		Start:      time.Now(),
		End:        time.Now().Add(1 * time.Hour),
		Step:       5,
		Period:     PeriodMinute,
		DaysInWeek: []int{1, 2, 3, 4, 5, 6, 7},
		Gap:        0,
	}
}

// parseTime converts various input types to time.Time
func parseTime(v interface{}) (time.Time, error) {
	switch val := v.(type) {
	case time.Time:
		return val, nil
	case string:
		s := strings.TrimSpace(val)
		
		// Try original string first (for RFC3339, ISO8601 with T)
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05.000Z",
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05",
		}
		for _, f := range formats {
			if t, err := time.Parse(f, s); err == nil {
				return t, nil
			}
		}
		
		// Try with space instead of T
		s = strings.Replace(s, "T", " ", 1)
		moreFormats := []string{
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02",
		}
		for _, f := range moreFormats {
			if t, err := time.Parse(f, s); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("unable to parse time string: %s", val)
	case int64:
		return time.Unix(val/1000, (val%1000)*1e6), nil // milliseconds
	case int:
		return time.Unix(int64(val)/1000, (int64(val)%1000)*1e6), nil
	default:
		return time.Time{}, fmt.Errorf("unsupported time type: %T", v)
	}
}

// toISOWeekday converts Go's Weekday (Sun=0) to ISO (Mon=1, Sun=7)
func toISOWeekday(wd time.Weekday) int {
	if wd == time.Sunday {
		return 7
	}
	return int(wd)
}

// GenerateSlots generates recurring time slots between start and end
func GenerateSlots(opts Options) ([]time.Time, error) {
	// Apply defaults
	if opts.Start == nil {
		opts.Start = DefaultOptions().Start
	}
	if opts.End == nil {
		opts.End = DefaultOptions().End
	}
	if opts.Step == 0 {
		opts.Step = DefaultOptions().Step
	}
	if opts.Period == "" {
		opts.Period = DefaultOptions().Period
	}
	if opts.DaysInWeek == nil {
		opts.DaysInWeek = DefaultOptions().DaysInWeek
	}

	start, err := parseTime(opts.Start)
	if err != nil {
		return nil, fmt.Errorf("invalid start: %w", err)
	}
	end, err := parseTime(opts.End)
	if err != nil {
		return nil, fmt.Errorf("invalid end: %w", err)
	}

	if start.After(end) {
		return []time.Time{}, nil
	}

	// Build day set for O(1) lookup
	daySet := make(map[int]bool)
	for _, d := range opts.DaysInWeek {
		daySet[d] = true
	}

	// Normalize to minute step
	stepMin := opts.Step
	switch opts.Period {
	case PeriodHour:
		stepMin = opts.Step*60 + opts.Gap
	case PeriodDay:
		stepMin = opts.Step * 1440
	case PeriodMinute:
		stepMin += opts.Gap
	}

	msPerStep := int64(stepMin) * 60_000 // milliseconds

	// Extract daily window from start/end
	startH, startM := start.Hour(), start.Minute()
	endH, endM := end.Hour(), end.Minute()
	crossMidnight := endH < startH || (endH == startH && endM <= startM)

	var slots []time.Time

	// Day granularity: simple loop
	if opts.Period == PeriodDay {
		for t := start.UnixMilli(); t <= end.UnixMilli(); t += msPerStep {
			d := time.UnixMilli(t)
			if daySet[toISOWeekday(d.Weekday())] {
				slots = append(slots, d)
			}
		}
		return slots, nil
	}

	// Minute/hour granularity: iterate days, then slots within window
	dayStart := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endMilli := end.UnixMilli()

	for day := dayStart.UnixMilli(); day <= endMilli; day += 86_400_000 {
		d := time.UnixMilli(day)
		dow := toISOWeekday(d.Weekday())
		if !daySet[dow] {
			continue
		}

		windowStart := day + int64(startH)*3_600_000 + int64(startM)*60_000
		windowEnd := day + int64(endH)*3_600_000 + int64(endM)*60_000
		if crossMidnight {
			windowEnd += 86_400_000 // next day
		}

		slotStart := windowStart
		if start.UnixMilli() > slotStart {
			slotStart = start.UnixMilli()
		}
		slotEnd := windowEnd
		if end.UnixMilli() < slotEnd {
			slotEnd = end.UnixMilli()
		}

		for t := slotStart; t < slotEnd; t += msPerStep {
			slots = append(slots, time.UnixMilli(t))
		}
	}

	return slots, nil
}
