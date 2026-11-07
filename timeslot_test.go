package timeslot

import (
	"testing"
	"time"
)


func TestGenerateSlots(t *testing.T) {
	tests := []struct {
		name      string
		opts      Options
		wantLen   int
		wantFirst time.Time
		wantLast  time.Time
	}{
		{
			name: "Test 1: Basic minute slots",
			opts: Options{
				Start:      "2022-03-12T05:30:00.000Z",
				End:        "2022-03-12T07:30:00.000Z",
				Step:       20,
				Period:     PeriodMinute,
				DaysInWeek: []int{6},
				Gap:        0,
			},
			wantLen:   6,
			wantFirst: time.Date(2022, 3, 12, 5, 30, 0, 0, time.UTC),
			wantLast:  time.Date(2022, 3, 12, 7, 10, 0, 0, time.UTC),
		},
		{
			name: "Test 2: Gap",
			opts: Options{
				Start:      "2022-03-11T05:30:00.000Z",
				End:        "2022-03-11T08:30:00.000Z",
				Step:       20,
				Period:     PeriodMinute,
				DaysInWeek: []int{5},
				Gap:        10,
			},
			wantLen:   6,
			wantFirst: time.Date(2022, 3, 11, 5, 30, 0, 0, time.UTC),
			wantLast:  time.Date(2022, 3, 11, 8, 0, 0, 0, time.UTC),
		},
		{
			name: "Test 3: Hour slots",
			opts: Options{
				Start:      "2022-03-12T05:30:00.000Z",
				End:        "2022-03-18T07:30:00.000Z",
				Step:       1,
				Period:     PeriodHour,
				DaysInWeek: []int{6, 1, 3},
				Gap:        0,
			},
			wantLen:   6,
			wantFirst: time.Date(2022, 3, 12, 5, 30, 0, 0, time.UTC),
			wantLast:  time.Date(2022, 3, 16, 6, 30, 0, 0, time.UTC),
		},
		{
			name: "Test 4: Single day",
			opts: Options{
				Start:      "2022-06-29T18:30:00.000Z",
				End:        "2022-06-29T20:30:00.000Z",
				Step:       20,
				Period:     PeriodMinute,
				DaysInWeek: []int{3, 4},
				Gap:        0,
			},
			wantLen:   6,
			wantFirst: time.Date(2022, 6, 29, 18, 30, 0, 0, time.UTC),
			wantLast:  time.Date(2022, 6, 29, 20, 10, 0, 0, time.UTC),
		},
		{
			name: "Test 5: Multi-week",
			opts: Options{
				Start:      "2022-06-29T05:30:00.000Z",
				End:        "2022-07-06T06:30:00.000Z",
				Step:       20,
				Period:     PeriodMinute,
				DaysInWeek: []int{6, 1, 3},
				Gap:        0,
			},
			wantLen:   12,
			wantFirst: time.Date(2022, 6, 29, 5, 30, 0, 0, time.UTC),
			wantLast:  time.Date(2022, 7, 6, 6, 10, 0, 0, time.UTC),
		},
		{
			name: "Edge: Cross-midnight",
			opts: Options{
				Start:      "2022-03-12T22:00:00.000Z",
				End:        "2022-03-13T06:00:00.000Z",
				Step:       1,
				Period:     PeriodHour,
				DaysInWeek: []int{6, 7},
				Gap:        0,
			},
			wantLen: 8,
		},
		{
			name: "Edge: Day period",
			opts: Options{
				Start:      "2022-03-14T00:00:00.000Z",
				End:        "2022-03-25T00:00:00.000Z",
				Step:       2,
				Period:     PeriodDay,
				DaysInWeek: []int{1, 2, 3, 4, 5},
				Gap:        0,
			},
			wantLen: 5,
		},
		{
			name: "Edge: Empty daysInWeek",
			opts: Options{
				Start:      "2022-03-12T05:30:00.000Z",
				End:        "2022-03-12T07:30:00.000Z",
				Step:       20,
				Period:     PeriodMinute,
				DaysInWeek: []int{},
				Gap:        0,
			},
			wantLen: 0,
		},
		{
			name: "Edge: Start after end",
			opts: Options{
				Start:      "2022-03-12T10:00:00.000Z",
				End:        "2022-03-12T09:00:00.000Z",
				Step:       30,
				Period:     PeriodMinute,
				DaysInWeek: []int{6},
				Gap:        0,
			},
			wantLen: 0,
		},
		{
			name: "Edge: Date object inputs",
			opts: Options{
				Start:      time.Date(2022, 3, 12, 5, 30, 0, 0, time.UTC),
				End:        time.Date(2022, 3, 12, 7, 30, 0, 0, time.UTC),
				Step:       60,
				Period:     PeriodMinute,
				DaysInWeek: []int{6},
				Gap:        0,
			},
			wantLen: 2,
		},
		{
			name: "Edge: Hour period with gap",
			opts: Options{
				Start:      "2022-03-12T05:30:00.000Z",
				End:        "2022-03-12T12:30:00.000Z",
				Step:       1,
				Period:     PeriodHour,
				DaysInWeek: []int{6},
				Gap:        30,
			},
			wantLen: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GenerateSlots(tc.opts)
			if err != nil {
				t.Fatalf("GenerateSlots error: %v", err)
			}
			if len(got) != tc.wantLen {
				t.Errorf("len=%d, want %d", len(got), tc.wantLen)
			}
			if !tc.wantFirst.IsZero() && len(got) > 0 {
				if !got[0].Equal(tc.wantFirst) {
					t.Errorf("first=%v, want %v", got[0], tc.wantFirst)
				}
			}
			if !tc.wantLast.IsZero() && len(got) > 0 {
				if !got[len(got)-1].Equal(tc.wantLast) {
					t.Errorf("last=%v, want %v", got[len(got)-1], tc.wantLast)
				}
			}
		})
	}
}
