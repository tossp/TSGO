package utils

import (
	"bytes"
	"strconv"
	"time"
)

func daysIn(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func Elapsed(from, to time.Time) (inverted bool, years, months, days, hours, minutes, seconds, nanoseconds int) {
	if from.Location() != to.Location() {
		to = to.In(to.Location())
	}

	inverted = false
	if from.After(to) {
		inverted = true
		from, to = to, from
	}

	y1, M1, d1 := from.Date()
	y2, M2, d2 := to.Date()

	h1, m1, s1 := from.Clock()
	h2, m2, s2 := to.Clock()

	ns1, ns2 := from.Nanosecond(), to.Nanosecond()

	years = y2 - y1
	months = int(M2 - M1)
	days = d2 - d1

	hours = h2 - h1
	minutes = m2 - m1
	seconds = s2 - s1
	nanoseconds = ns2 - ns1

	if nanoseconds < 0 {
		nanoseconds += 1e9
		seconds--
	}
	if seconds < 0 {
		seconds += 60
		minutes--
	}
	if minutes < 0 {
		minutes += 60
		hours--
	}
	if hours < 0 {
		hours += 24
		days--
	}
	if days < 0 {
		days += daysIn(y2, M2-1)
		months--
	}
	if months < 0 {
		months += 12
		years--
	}
	return
}

func HumanElapsed(from, to time.Time) string {
	_, years, months, days, hours, minutes, seconds, nanoseconds := Elapsed(from, to)
	var b bytes.Buffer
	if years > 0 {
		b.WriteString(strconv.Itoa(years))
		b.WriteString("年")
	}
	if months > 0 {
		b.WriteString(strconv.Itoa(months))
		b.WriteString("月")
	}
	if days > 0 {
		b.WriteString(strconv.Itoa(days))
		b.WriteString("天")
	}
	if hours > 0 {
		b.WriteString(strconv.Itoa(hours))
		b.WriteString("时")
	}
	if minutes > 0 {
		b.WriteString(strconv.Itoa(minutes))
		b.WriteString("分")
	}
	if seconds > 0 {
		b.WriteString(strconv.Itoa(seconds))
		b.WriteString("秒")
	}
	if nanoseconds > 0 {
		b.WriteString(strconv.Itoa(nanoseconds))
		b.WriteString("纳秒")
	}
	return b.String()
}
