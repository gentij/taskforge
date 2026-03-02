package utils

import (
	"strconv"
	"time"
)

func RelativeTime(now time.Time, t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	d := now.Sub(t)
	if d < 0 {
		d = -d
	}
	seconds := int(d.Seconds())
	if seconds < 60 {
		return "just now"
	}
	minutes := seconds / 60
	if minutes < 60 {
		return formatAge(minutes, "m")
	}
	hours := minutes / 60
	if hours < 24 {
		return formatAge(hours, "h")
	}
	days := hours / 24
	if days < 7 {
		return formatAge(days, "d")
	}
	weeks := days / 7
	return formatAge(weeks, "w")
}

func formatAge(value int, unit string) string {
	return formatInt(value) + unit + " ago"
}

func formatInt(value int) string {
	return strconv.Itoa(value)
}
