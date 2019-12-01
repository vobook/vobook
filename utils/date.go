package utils

import (
	"math"
	"time"
)

type DateDiffResult struct {
	Year,
	Month,
	Day,
	Hour,
	Min,
	Sec int
}

// AgeAt gets the age of an entity at a certain time.
func AgeAt(birthDate time.Time, now time.Time) int {
	years := now.Year() - birthDate.Year()

	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Age(birthDate time.Time) int {
	return AgeAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}

func BirthdayInfo(bornAt time.Time) (nextAt time.Time, daysLeft, currentAge int) {
	_, mo1, d1 := bornAt.Date()
	y2, mo2, _ := time.Now().Date()

	// adjust year if birthday passed
	if mo2 > mo1 {
		y2++
	}
	nextAt = time.Date(y2, mo1, d1, 0, 0, 0, 0, bornAt.Location())
	daysLeft = DaysBetween(nextAt, time.Now())
	currentAge = Age(bornAt)

	return
}

func DaysBetween(a, b time.Time) int {
	if a.After(b) {
		a, b = b, a
	}

	if b.Sub(a).Hours()/24.0 < 1 {
		return 0
	}

	return int(math.Ceil(b.Sub(a).Hours() / 24.0))
}

func DateDiff(a, b time.Time) DateDiffResult {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}

	y1, mo1, d1 := a.Date()
	y2, mo2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year := y2 - y1
	month := int(mo2 - mo1)
	day := d2 - d1
	hour := h2 - h1
	min := m2 - m1
	sec := s2 - s1

	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, mo1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return DateDiffResult{
		Year:  year,
		Month: month,
		Day:   day,
		Hour:  hour,
		Min:   min,
		Sec:   sec,
	}
}
