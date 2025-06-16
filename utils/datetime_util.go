package utils

import (
	models "LeafMS-BackEnd/models"
	"errors"
	"log"
	"sort"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func rollLeaveBackwardOneDay(date models.Datetime) models.Datetime {
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

func rollLeaveForwardOneDay(date models.Datetime) models.Datetime {
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

func rollLeaveBackward(date models.Datetime, daysBackward int) models.Datetime {
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

func rollLeaveForward(date models.Datetime, daysForward int) models.Datetime {
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

// ================================
// Main Functions
// ================================
func RemoveWeekendsFromLeaveData(leavesSpan models.LeaveInfo) ([]models.LeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	leaveStartDate, err := ParseStringToDate(leavesSpan.StartDate)
	if err != nil {
		return nil, err
	}
	leaveEndDate, err := ParseStringToDate(leavesSpan.EndDate)
	if err != nil {
		return nil, err
	}

	currentDate := leaveStartDate

	weekdayInInt, _ := DateToWeekday(currentDate)
	if weekdayInInt == 0 {
		currentDate = rollLeaveForward(currentDate, 1)
		weekdayInInt = 1
	} else if weekdayInInt == 6 {
		currentDate = rollLeaveForward(currentDate, 2)
		weekdayInInt = 1
	}

	for leaveEndDate.IsGreaterThanOrEquals(currentDate) {
		var leaveSpan models.LeaveInfo
		if endDate := rollLeaveForward(currentDate, 5-weekdayInInt); leaveEndDate.IsGreaterThanOrEquals(endDate) {
			leaveSpan = models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: ParseDateToString(currentDate),
				EndDate:   ParseDateToString(endDate),
				Reason:    leavesSpan.Reason,
			}
		} else {
			leaveSpan = models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: ParseDateToString(currentDate),
				EndDate:   ParseDateToString(leaveEndDate),
				Reason:    leavesSpan.Reason,
			}
		}

		splitLeaves = append(splitLeaves, leaveSpan)
		currentDate = rollLeaveForward(currentDate, 8-weekdayInInt)
		weekdayInInt = 1
	}

	return splitLeaves, nil
}

func FetchHolidays(leave models.LeaveInfo) ([]models.Holiday, error) {
	leaveStartDate, err := ParseStringToDate(leave.StartDate)
	if err != nil {
		log.Println("There was problem parsing the starting date of a leave Err:", err)
		return []models.Holiday{}, err
	}
	if err = FeasibleDate(leaveStartDate); err != nil {
		log.Println("The start date is not practically possible in the real world. Err: ", err)
		return []models.Holiday{}, err
	}

	leaveEndDate, err := ParseStringToDate(leave.EndDate)
	if err != nil {
		log.Println("There was problem parsing the ending date of a leave Err:", err)
		return []models.Holiday{}, err
	}
	if err = FeasibleDate(leaveEndDate); err != nil {
		log.Println("The start date is not practically possible in the real world. Err: ", err)
		return []models.Holiday{}, err
	}

	holidaysBson, err := database.Find("publicHolidays", bson.D{
		{Key: "$and", Value: bson.A{
			bson.D{{Key: "date.datetime.year", Value: bson.D{
				{Key: "$gte", Value: leaveStartDate.Year},
				{Key: "$lte", Value: leaveEndDate.Year},
			}}},
			bson.D{{Key: "date.datetime.month", Value: bson.D{
				{Key: "$gte", Value: leaveStartDate.Month},
				{Key: "$lte", Value: leaveEndDate.Month},
			}}},
			bson.D{{Key: "date.datetime.day", Value: bson.D{
				{Key: "$gte", Value: leaveStartDate.Day},
				{Key: "$lte", Value: leaveEndDate.Day},
			}}},
		}},
	})
	if err != nil {
		errMessage := "For fuck's sake, there was a problem, "
		errMessage += "while trying to find a holiday conflicting with the applied leave in the database. Err:"
		log.Println(errMessage, err)
		return []models.Holiday{}, err
	}
	holidays := ReturnHolidays(holidaysBson)
	sort.Sort(Holidays(holidays))

	return holidays, nil
}

func RemoveHolidayFromLeaveData(leave models.LeaveInfo) ([]models.LeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	holidays, err := FetchHolidays(leave)
	if err != nil {
		log.Println(err)
	}

	startDate := leave.StartDate
	for _, holiday := range holidays {
		parsedStartDate, _ := ParseStringToDate(startDate)
		if !parsedStartDate.IsGreaterThanOrEquals(holiday.Date.Datetime) {
			leaveSpan := models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: startDate,
				EndDate:   ParseDateToString(rollLeaveBackwardOneDay(holiday.Date.Datetime)),
			}
			splitLeaves = append(splitLeaves, leaveSpan)
		}
		startDate = ParseDateToString(rollLeaveForwardOneDay(holiday.Date.Datetime))
	}

	parsedStartDate, _ := ParseStringToDate(startDate)
	parsedEndDate, _ := ParseStringToDate(leave.EndDate)
	if !parsedStartDate.IsGreaterThanOrEquals(parsedEndDate) {
		lastLeaveSpan := models.LeaveInfo{
			Id:        primitive.NewObjectID(),
			StartDate: startDate,
			EndDate:   leave.EndDate,
		}
		splitLeaves = append(splitLeaves, lastLeaveSpan)
	}
	return splitLeaves, nil
}

type Holidays []models.Holiday

func (holidays Holidays) Len() int      { return len(holidays) }
func (holidays Holidays) Swap(i, j int) { holidays[i], holidays[j] = holidays[j], holidays[i] }
func (holidays Holidays) Less(i, j int) bool {
	return !holidays[i].Date.Datetime.IsGreaterThanOrEquals(holidays[j].Date.Datetime)
}
