package cron

import (
	"errors"
	"strconv"
	"strings"
)

// Schedule represents a cron schedule.
type Schedule struct {
	expr     string // formatted expression
	minutes  set64  // values: 0-59
	hours    set64  // values: 0-23
	days     set64  // values: 1-31
	months   set64  // values: 1-12
	weekDays set64  // values: 0-6
}

// Parse a cron schedule.
//
// See https://wikipedia.org/wiki/Cron
func ParseSchedule(s string) (Schedule, error) {
	if macros[s] != "" {
		s = macros[s]
	}

	ss := strings.Fields(s)
	if len(ss) != 5 {
		return Schedule{}, errBadFormat
	}

	sch := Schedule{
		expr: s,
	}

	var ok [5]bool
	sch.minutes, ok[0] = parseField(ss[0], 0, 59, nil)
	sch.hours, ok[1] = parseField(ss[1], 0, 23, nil)
	sch.days, ok[2] = parseField(ss[2], 1, 31, nil)
	sch.months, ok[3] = parseField(ss[3], 1, 12, monthNames)
	sch.weekDays, ok[4] = parseField(ss[4], 0, 6, dayNames)

	if ok != [5]bool{true, true, true, true, true} {
		return Schedule{}, errInvalidSchedule
	}
	return sch, nil
}

func parseField(s string, min, max int, aliases map[string]int) (set64, bool) {
	var m set64
	for _, s := range strings.Split(s, ",") {
		a, b := s, s
		if i := strings.IndexByte(s, '-'); i >= 0 {
			a, b = s[:i], s[i+1:]
		}

		lo := parseToken(a, min, aliases)
		hi := parseToken(b, max, aliases)

		if lo < min || max < hi || hi < lo {
			return 0, false
		}

		for i := lo; i <= hi; i++ {
			m.set(i)
		}
	}
	return m, true
}

func parseToken(s string, wild int, aliases map[string]int) int {
	if s == "*" {
		return wild
	}
	if n, ok := aliases[strings.ToUpper(s)]; ok {
		return n
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return -1
}

func (s Schedule) String() string {
	return s.expr
}

// set64 is a set of integers for [0..63] range.
type set64 uint64

func (s *set64) set(i int)     { *s |= 1 << uint(i) }
func (s set64) has(i int) bool { return s&(1<<uint(i)) != 0 }

var (
	errBadFormat       = errors.New("bad format")
	errInvalidSchedule = errors.New("invalid schedule")

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
