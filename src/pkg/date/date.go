package date

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)
var isLastDayOfWeekSaturday = true
func LastDayOfWeekIsSaturday(b bool) {
	isLastDayOfWeekSaturday = b
}

func WeekStart(year, week int) time.Time {
	T := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	if wd := T.Weekday(); wd == time.Sunday {
		T = T.AddDate(0, 0, -6)
	} else {
		T = T.AddDate(0, 0, -int(wd)+1)
	}
	_, w := T.ISOWeek()
	T = T.AddDate(0, 0, (week-w)*7)
	return T
}

func WeekRange(year, week int) (start, end time.Time) {
    start = WeekStart(year, week)
    end = start.AddDate(0, 0, 6)
	if isLastDayOfWeekSaturday {
		end = end.Add(-1)
	}
    return
}

func WeekDay(i interface{}) (wkd time.Weekday,err error) {
	var s string
	var ok bool
	if reflect.ValueOf(i).Kind() == reflect.Interface {
		s, ok = i.(string)
		if !ok {
			return wkd, fmt.Errorf("day is invalid")
		}
	}
	s, _ = i.(string)
	switch strings.ToLower(s) {
	case "monday", "mon":
		wkd = time.Monday
	case "tuesday", "tue":
		wkd = time.Tuesday
	case "wednesday", "wed":
		wkd = time.Wednesday
	case "thursday", "thur":
		wkd = time.Thursday
	case "friday", "fri":
		wkd = time.Friday
	case "saturday", "sat":
		wkd = time.Saturday
	case "sunday", "sun":
		wkd = time.Sunday
	default:
		return wkd, fmt.Errorf("weekday is invalid")
	}
	return wkd, nil
}