package main

import (
	"testing"
	"time"
)

func TestBetweenTime(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->23:59", now)
	if got != true {
		t.Errorf("Should be within time")
	}
}

func TestOutsideTime(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->15:00", now)
	if got == true {
		t.Errorf("Should be whithout time")
	}
}

func TestOutsideTimeInsideMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("00:00->15:00;January", now)
	if got != true {
		t.Errorf("Should be whithin time")
	}
}

func TestDayOfTheWeek(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("Wednesday", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestOutsideDayOfTheWeek(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("Thursday", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}

func TestDayAndMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January 01", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestOutsideDayAndMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January 02", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("January", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonthMultipleLastCorrect(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February;January", now)
	if got != true {
		t.Errorf("Should be correct for date")
	}
}

func TestMonthMultipleFirstCorrect(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February;January", now)
	if got != true {
		t.Errorf("Should be correct for Month")
	}
}

func TestOutsideMonth(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February", now)
	if got == true {
		t.Errorf("Should be incorrect for Month")
	}
}

func TestInsideCustomWeekDay(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2020-01-01T15:01:00")
	got := isWithinTime("February", now)
	if got == true {
		t.Errorf("Should be correct for date")
	}
}
