package util

import (
	"strings"
	"time"

	"github.com/oapi-codegen/runtime/types"
)

type timeBlock struct{}

func (t timeBlock) formatter(rtn *[]string, tm time.Time, r rune, zeroLeading bool) {
	add := func(value string) {
		*rtn = append(*rtn, value)
	}

	var addInt func(int)
	if zeroLeading {
		addInt = func(value int) {
			if value < 10 {
				add("0")
			}
			add(String.FromInt(value))
		}
	} else {
		addInt = func(value int) {
			add(String.FromInt(value))
		}
	}

	var iVal int
	switch r {
	case 'd':
		addInt(tm.Day())
	case 'F':
		t.formatter(rtn, tm, 'Y', zeroLeading)
		add("-")
		t.formatter(rtn, tm, 'm', zeroLeading)
		add("-")
		t.formatter(rtn, tm, 'd', zeroLeading)
	case 'H':
		addInt(tm.Hour())
	case 'I':
		iVal = tm.Hour() % 12
		if iVal == 0 {
			iVal = 12
		}
		addInt(iVal)
	case 'm':
		addInt(int(tm.Month()))
	case 'M':
		addInt(tm.Minute())
	case 'R':
		t.formatter(rtn, tm, 'H', zeroLeading)
		add(":")
		t.formatter(rtn, tm, 'M', zeroLeading)
	case 'S':
		addInt(tm.Second())
	case 'T':
		t.formatter(rtn, tm, 'H', zeroLeading)
		add(":")
		t.formatter(rtn, tm, 'M', zeroLeading)
		add(":")
		t.formatter(rtn, tm, 'S', zeroLeading)
	case 'w':
		addInt(int(tm.Weekday()))
	case 'Y':
		addInt(tm.Year())
	case '%':
		add("%")
	}
}

func (t timeBlock) parserFormat(format string) string {
	format = strings.Replace(format, "%F", "%Y-%m-%d", -1)
	format = strings.Replace(format, "%T", "%R:%S", -1)
	format = strings.Replace(format, "%R", "%H:%M", -1)
	format = strings.Replace(format, "%Y", "2006", -1)
	format = strings.Replace(format, "%m", "01", -1)
	format = strings.Replace(format, "%d", "02", -1)
	format = strings.Replace(format, "%H", "15", -1)
	format = strings.Replace(format, "%M", "04", -1)
	format = strings.Replace(format, "%S", "05", -1)
	format = strings.Replace(format, "%a", "Mon", -1)
	format = strings.Replace(format, "%b", "Jan", -1)
	return format
}

// Sprintf : like strftime
func (t timeBlock) Sprintf(format string, tm time.Time) string {
	var rtn []string

	runes := []rune(format)

	spec := false
	zeroLeading := false
	for _, r := range runes {
		if spec {
			if r == '0' {
				zeroLeading = true
			} else {
				t.formatter(&rtn, tm, r, zeroLeading)
				spec = false
			}
		} else {
			if r == '%' {
				spec = true
			} else {
				rtn = append(rtn, string(r))
			}
		}
	}

	return String.Concat(rtn)
}

// Parse :
func (t timeBlock) Parse(format string, value string) (time.Time, error) {
	return time.Parse(t.parserFormat(format), value)
}

// ParseWithTz :
func (t timeBlock) ParseWithTz(format, value, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation(t.parserFormat(format), value, loc)
}

// UnixMilli :
func (t timeBlock) UnixMilli(tm time.Time) int64 {
	if tm.IsZero() {
		return -1
	}

	return tm.UnixNano() / 1000000
}

// ToDate :
func (t timeBlock) ToDate(tm time.Time, bias int) string {
	return t.Sprintf("%0F", tm.AddDate(0, 0, bias))
}

// ToOapiDate :
func (t timeBlock) ToOapiDate(tm time.Time) types.Date {
	return types.Date{Time: tm}
}

// ToTimeFromOapiDate :
func (t timeBlock) ToTimeFromOapiDate(date types.Date) time.Time {
	if t.OapiDateNilCheck(date) {
		return time.Time{}
	}

	parsedTime, err := time.Parse(time.DateOnly, date.Format(time.DateOnly))
	if err != nil {
		return time.Time{}
	}

	return parsedTime
}

// OapiDateNilCheck :
func (t timeBlock) OapiDateNilCheck(date types.Date) bool {
	return date.Time.IsZero()
}
