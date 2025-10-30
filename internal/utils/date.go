package utils

import "time"

func TruncateToMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func MonthsOverlap(aStart, aEnd, bStart, bEnd time.Time) int {
	if aStart.After(aEnd) {
		return 0
	}
	if bStart.After(bEnd) {
		return 0
	}

	start := maxTime(aStart, bStart)
	end := minTime(aEnd, bEnd)

	if start.After(end) {
		return 0
	}

	return monthsBeforeInclusive(start, end)

}

func monthsBeforeInclusive(a, b time.Time) int {
	years := b.Year() - a.Year()

	months := int(b.Month()) - int(a.Month())

	return years*12 + months + 1
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
