package cron

// Schedule represents a cron schedule.
type Schedule struct {
	expr     string // formatted expression
	minutes  set64  // values: 0-59
	hours    set64  // values: 0-23
	days     set64  // values: 1-31
	months   set64  // values: 1-12
	weekDays set64  // values: 0-6
}

func (s Schedule) String() string {
	return s.expr
}

// set64 is a set of integers for [0..63] range.
type set64 uint64

func (s *set64) set(i int)     { *s |= 1 << uint(i) }
func (s set64) has(i int) bool { return s&(1<<uint(i)) != 0 }

var (
	macros = map[string]string{
		"@yearly":   "0 0 1 1 *",
		"@annually": "0 0 1 1 *",
		"@monthly":  "0 0 1 * *",
		"@weekly":   "0 0 * * 0",
		"@daily":    "0 0 * * *",
		"@midnight": "0 0 * * *",
		"@hourly":   "0 * * * *",
		// TODO: @reboot per service start?
	}

	monthNames = map[string]int{
		"JAN": 1,
		"FEB": 2,
		"MAR": 3,
		"APR": 4,
		"MAY": 5,
		"JUN": 6,
		"JUL": 7,
		"AUG": 8,
		"SEP": 9,
		"OCT": 10,
		"NOV": 11,
		"DEC": 12,
	}

	dayNames = map[string]int{
		"SUN": 0,
		"MON": 1,
		"TUE": 2,
		"WED": 3,
		"THU": 4,
		"FRI": 5,
		"SAT": 6,
	}
)
