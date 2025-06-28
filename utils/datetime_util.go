package utils

import (
	models "LeafMS-BackEnd/models"
	"errors"
)

var daysInMonth = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}
var WeekDays = map[int]string{
	0: "Sunday",
	1: "Monday",
	2: "Tuesday",
	3: "Wednesday",
	4: "Thursday",
	5: "Friday",
	6: "Saturday",
}

// ================================
// Helper Functions
// ================================
func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	} else if year%4 == 0 && year%100 != 0 {
		return true
	}
	return false
}

func FeasibleDate(date models.Datetime) error {
	day := date.Day
	month := date.Month
	year := date.Year

	if month > 12 {
		errStr := "month provided is more than what is practically possible.\n"
		errStr += "i.e - The month is 13th or more than 13th, which is not possible"
		err := errors.New(errStr)
		return err
	} else if (isLeapYear(year) && month == 2 && day > (daysInMonth[month]+1)) || (!isLeapYear(year) && day > daysInMonth[month]) {
		errStr := "The number of days is more than possible for the month in the date.\n"
		err := errors.New(errStr)
		return err
	}
	return nil
}

func DateToWeekday(date models.Datetime) (int, string) {
	year := date.Year - 2000
	year += (year / 4) + 7

	for i := 1; i < date.Month; i++ {
		year += daysInMonth[i]
	}
	year += date.Day - 1
	if isLeapYear(date.Year) && date.Month <= 2 {
		year -= 1
	}
	year %= 7
	return year, WeekDays[year]
}

func RollLeaveBackwardOneDay(date models.Datetime) models.Datetime {
	if date.Day == 1 {
		if date.Month == 1 {
			date.Year -= 1
			date.Month = 12
		} else {
			date.Month -= 1
		}
		date.Day = daysInMonth[date.Month]
		if date.Month == 2 && isLeapYear(date.Year) {
			date.Day += 1
		}
	} else {
		date.Day -= 1
	}
	return date
}

func RollLeaveForwardOneDay(date models.Datetime) models.Datetime {
	if date.Day == daysInMonth[date.Month] ||
		(date.Month == 2 && isLeapYear(date.Year) && date.Day == daysInMonth[date.Month]+1) {
		date.Day = 1
		if date.Month == 12 {
			date.Year += 1
			date.Month = 1
		} else {
			date.Month += 1
		}
	} else {
		date.Day += 1
	}
	return date
}

func RollLeaveBackward(date models.Datetime, daysBackward int) models.Datetime {
	if date.Day-daysBackward <= 0 {
		if date.Month == 1 {
			date.Year -= 1
			date.Month = 12
		} else {
			date.Month -= 1
		}
		date.Day = daysInMonth[date.Month] - (daysBackward - date.Day)
		if date.Month == 2 && isLeapYear(date.Year) {
			date.Day += 1
		}
	} else {
		date.Day -= daysBackward
	}
	return date
}

func RollLeaveForward(date models.Datetime, daysForward int) models.Datetime {
	if (date.Day+daysForward) > daysInMonth[date.Month] ||
		(date.Month == 2 && isLeapYear(date.Year) && (date.Day+daysForward) > daysInMonth[date.Month]+1) {
		date.Day = (date.Day + daysForward) - daysInMonth[date.Month]
		if date.Month == 2 && isLeapYear(date.Year) && (date.Day+daysForward) > daysInMonth[date.Month]+1 {
			date.Day -= 1
		}
		if date.Month == 12 {
			date.Year += 1
			date.Month = 1
		} else {
			date.Month += 1
		}
	} else {
		date.Day += daysForward
	}
	return date
}
