package repo

import (
	"testing"
	"time"
)

func Test_fillEmptyDays(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour)
	startDate := now.AddDate(0, 0, -29)

	tests := []struct {
		name     string
		input    []DailyCount
		expected int // expected number of days (always 30)
		nonZero  map[time.Time]int64
	}{
		{
			name:     "no input days",
			input:    []DailyCount{},
			expected: 30,
			nonZero:  map[time.Time]int64{},
		},
		{
			name: "some days filled",
			input: []DailyCount{
				{Day: startDate.AddDate(0, 0, 0), Count: 5},
				{Day: startDate.AddDate(0, 0, 10), Count: 10},
				{Day: startDate.AddDate(0, 0, 29), Count: 2},
			},
			expected: 30,
			nonZero: map[time.Time]int64{
				startDate.AddDate(0, 0, 0):  5,
				startDate.AddDate(0, 0, 10): 10,
				startDate.AddDate(0, 0, 29): 2,
			},
		},
		{
			name: "all days filled",
			input: func() []DailyCount {
				var res []DailyCount
				for i := 0; i < 30; i++ {
					day := startDate.AddDate(0, 0, i)
					res = append(res, DailyCount{Day: day, Count: int64(i)})
				}
				return res
			}(),
			expected: 30,
			nonZero: func() map[time.Time]int64 {
				m := make(map[time.Time]int64)
				for i := 0; i < 30; i++ {
					day := startDate.AddDate(0, 0, i)
					m[day] = int64(i)
				}
				return m
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fillEmptyDays(tt.input)

			if len(got) != tt.expected {
				t.Errorf("expected %d days, got %d", tt.expected, len(got))
			}

			for _, dc := range got {
				wantCount := tt.nonZero[dc.Day]
				if dc.Count != wantCount {
					t.Errorf("on day %s: expected count %d, got %d", dc.Day.Format("2006-01-02"), wantCount, dc.Count)
				}
			}
		})
	}
}
