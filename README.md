# timeslot

Go module for generating recurring time slots from a date-time range with flexible configuration options.

## Features

- **Flexible time input**: Accepts `time.Time`, RFC3339/ISO8601 strings, and Unix timestamps (milliseconds)
- **Multiple periods**: Minute, Hour, and Day granularity
- **Weekday filtering**: Generate slots only on specific days of the week (ISO 1-7, Mon-Sun)
- **Configurable step & gap**: Set interval between slots and optional gap
- **Cross-midnight support**: Handles ranges that span across midnight
- **Zero dependencies**: Pure Go standard library

## Installation

```bash
go get github.com/smdpro/timeslot
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/your-org/timeslot"
)

func main() {
	// Generate 15-minute slots for business hours (9am-5pm), Mon-Fri, for one week
	slots, err := timeslot.GenerateSlots(timeslot.Options{
		Start:      "2024-01-01T09:00:00",
		End:        "2024-01-07T17:00:00",
		Step:       15,
		Period:     timeslot.PeriodMinute,
		DaysInWeek: []int{1, 2, 3, 4, 5}, // Mon-Fri
	})
	if err != nil {
		panic(err)
	}

	for _, slot := range slots {
		fmt.Println(slot.Format("2006-01-02 15:04"))
	}
}
```

## Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Start` | `interface{}` | `time.Now()` | Start of range. Accepts `time.Time`, string (RFC3339/ISO8601), or int64 (ms timestamp) |
| `End` | `interface{}` | `time.Now() + 1h` | End of range. Same types as Start |
| `Step` | `int` | `5` | Interval between slots |
| `Period` | `Period` | `PeriodMinute` | Unit for Step: `PeriodMinute` ("m"), `PeriodHour` ("h"), `PeriodDay` ("d") |
| `DaysInWeek` | `[]int` | `[1,2,3,4,5,6,7]` | ISO weekdays to include (1=Mon, 7=Sun) |
| `Gap` | `int` | `0` | Additional gap in minutes between slots (for Minute/Hour periods) |

### Period Constants

```go
const (
    PeriodMinute Period = "m"
    PeriodHour   Period = "h"
    PeriodDay    Period = "d"
)
```

## Supported Time Formats

The `Start` and `End` fields accept:

| Type | Examples |
|------|----------|
| `time.Time` | `time.Now()`, `time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)` |
| String (RFC3339) | `"2024-01-01T09:00:00Z"`, `"2024-01-01T09:00:00.000Z"` |
| String (ISO8601) | `"2024-01-01T09:00:00"`, `"2024-01-01T09:00"` |
| String (space separator) | `"2024-01-01 09:00:00"`, `"2024-01-01 09:00"`, `"2024-01-01"` |
| int64 (ms timestamp) | `1704099600000` |
| int (ms timestamp) | `1704099600000` |

## Examples

### 1. Hourly slots for a day
```go
slots, _ := timeslot.GenerateSlots(timeslot.Options{
    Start:  "2024-01-01T00:00:00",
    End:    "2024-01-01T23:00:00",
    Step:   1,
    Period: timeslot.PeriodHour,
})
// Returns: 2024-01-01 00:00, 01:00, 02:00, ..., 23:00
```

### 2. Daily slots for weekdays only
```go
slots, _ := timeslot.GenerateSlots(timeslot.Options{
    Start:      "2024-01-01",
    End:        "2024-01-31",
    Step:       1,
    Period:     timeslot.PeriodDay,
    DaysInWeek: []int{1, 2, 3, 4, 5}, // Mon-Fri
})
// Returns: All weekdays in January 2024 at 00:00
```

### 3. 30-min slots with 15-min gap (e.g., appointment buffers)
```go
slots, _ := timeslot.GenerateSlots(timeslot.Options{
    Start:  "2024-01-01T09:00:00",
    End:    "2024-01-01T17:00:00",
    Step:   30,
    Period: timeslot.PeriodMinute,
    Gap:    15, // 15-min buffer between appointments
})
// Returns: 09:00, 09:45, 10:30, 11:15, ... (30 min slot + 15 min gap = 45 min step)
```

### 4. Cross-midnight range (e.g., night shift)
```go
slots, _ := timeslot.GenerateSlots(timeslot.Options{
    Start:  "2024-01-01T22:00:00",
    End:    "2024-01-02T06:00:00",
    Step:   60,
    Period: timeslot.PeriodMinute,
})
// Returns: 22:00, 23:00, 00:00, 01:00, ..., 06:00 (spans across midnight)
```

### 5. Using time.Time directly
```go
loc, _ := time.LoadLocation("America/New_York")
slots, _ := timeslot.GenerateSlots(timeslot.Options{
    Start: time.Date(2024, 1, 1, 9, 0, 0, 0, loc),
    End:   time.Date(2024, 1, 1, 17, 0, 0, 0, loc),
    Step:  30,
    Period: timeslot.PeriodMinute,
})
```

## Default Options

```go
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
```

## Error Handling

`GenerateSlots` returns an error if:
- `Start` or `End` cannot be parsed (unsupported type or invalid format)
- Internal time conversion fails

```go
slots, err := timeslot.GenerateSlots(timeslot.Options{
    Start: "invalid-date",
    End:   "2024-01-01",
})
if err != nil {
    // err: "invalid start: unable to parse time string: invalid-date"
}
```

## Edge Cases

- **Start after End**: Returns empty slice `[]time.Time{}`, no error
- **Empty DaysInWeek**: Uses default (all days)
- **Zero Step**: Uses default (5 minutes)
- **Empty Period**: Uses default (Minute)
- **Timezone**: Preserved from input `time.Time`; string inputs parsed as UTC unless offset specified

## License

MIT
