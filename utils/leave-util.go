package utils

import (
	"LeafMS-BackEnd/database"
	models "LeafMS-BackEnd/models"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RemoveWeekendsFromLeaveData(leavesSpan models.LeaveInfo) ([]models.LeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	leaveStartDate := models.Datetime{Year: leavesSpan.StartDate.Year(), Month: int(leavesSpan.StartDate.Month()), Day: leavesSpan.StartDate.Day()}
	leaveEndDate := models.Datetime{Year: leavesSpan.EndDate.Year(), Month: int(leavesSpan.EndDate.Month()), Day: leavesSpan.EndDate.Day()}

	currentDate := leaveStartDate

	weekdayInInt, _ := DateToWeekday(currentDate)
	if weekdayInInt == 0 {
		currentDate = RollLeaveForward(currentDate, 1)
		weekdayInInt = 1
	} else if weekdayInInt == 6 {
		currentDate = RollLeaveForward(currentDate, 2)
		weekdayInInt = 1
	}

	for leaveEndDate.IsGreaterThanOrEquals(currentDate) {
		var leaveSpan models.LeaveInfo
		if endDate := RollLeaveForward(currentDate, 5-weekdayInInt); leaveEndDate.IsGreaterThanOrEquals(endDate) {
			leaveSpan = models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: time.Date(currentDate.Year, time.Month(currentDate.Month), currentDate.Day, 0, 0, 0, 0, time.Local),
				EndDate:   time.Date(endDate.Year, time.Month(endDate.Month), endDate.Day, 0, 0, 0, 0, time.Local),
				Reason:    leavesSpan.Reason,
			}
		} else {
			leaveSpan = models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: time.Date(currentDate.Year, time.Month(currentDate.Month), currentDate.Day, 0, 0, 0, 0, time.Local),
				EndDate:   time.Date(leaveEndDate.Year, time.Month(leaveEndDate.Month), leaveEndDate.Day, 0, 0, 0, 0, time.Local),
				Reason:    leavesSpan.Reason,
			}
		}

		splitLeaves = append(splitLeaves, leaveSpan)
		currentDate = RollLeaveForward(currentDate, 8-weekdayInInt)
		weekdayInInt = 1
	}

	return splitLeaves, nil
}

func FetchHolidaysBetweenRequestedLeave(leave models.LeaveInfo) ([]models.Holiday, error) {
	leaveStartDate := models.Datetime{Year: leave.StartDate.Year(), Month: int(leave.StartDate.Month()), Day: leave.StartDate.Day()}

	if err := FeasibleDate(leaveStartDate); err != nil {
		log.Println("The start date is not practically possible in the real world. Err: ", err)
		return []models.Holiday{}, err
	}

	leaveEndDate := models.Datetime{Year: leave.EndDate.Year(), Month: int(leave.EndDate.Month()), Day: leave.EndDate.Day()}
	if err := FeasibleDate(leaveEndDate); err != nil {
		log.Println("The start date is not practically possible in the real world. Err: ", err)
		return []models.Holiday{}, err
	}

	holidaysBson, err := DbConn.Find("publicHolidays", bson.D{
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
	holidays := database.ConvertRawBsonToHolidays(holidaysBson)
	sort.Sort(Holidays(holidays))

	return holidays, nil
}

func RemoveHolidayFromLeaveData(grossLeave models.LeaveInfo) ([]models.LeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	holidays, err := FetchHolidaysBetweenRequestedLeave(grossLeave)
	if err != nil {
		log.Println(err)
	}

	startDate := models.Datetime{Year: grossLeave.StartDate.Year(), Month: int(grossLeave.StartDate.Month()), Day: grossLeave.StartDate.Day()}
	for _, holiday := range holidays {
		if !startDate.IsGreaterThanOrEquals(holiday.Date.Datetime) {
			rolledBackLeave := RollLeaveBackwardOneDay(holiday.Date.Datetime)
			leaveSpan := models.LeaveInfo{
				Id:        primitive.NewObjectID(),
				StartDate: time.Date(startDate.Year, time.Month(startDate.Month), startDate.Day, 0, 0, 0, 0, time.Local),
				EndDate:   time.Date(rolledBackLeave.Year, time.Month(rolledBackLeave.Month), rolledBackLeave.Day, 0, 0, 0, 0, time.Local),
			}
			splitLeaves = append(splitLeaves, leaveSpan)
		}
		startDate = RollLeaveForwardOneDay(holiday.Date.Datetime)
	}

	grossEndDate := models.Datetime{Year: grossLeave.StartDate.Year(), Month: int(grossLeave.StartDate.Month()), Day: grossLeave.StartDate.Day()}
	if !startDate.IsGreaterThanOrEquals(grossEndDate) {
		lastLeaveSpan := models.LeaveInfo{
			Id:        primitive.NewObjectID(),
			StartDate: time.Date(startDate.Year, time.Month(startDate.Month), startDate.Day, 0, 0, 0, 0, time.Local),
			EndDate:   grossLeave.EndDate,
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
