package utils

import (
	"log"
	"strconv"

	models "LeafMS-BackEnd/models"
)

func ParseStringToDate(date string) (models.Datetime, error) {
	const (
		Date = iota
		Month
	)

	var valToInsert = Date
	var parsedDate models.Datetime
	var currParserVal = ""
	var index = 0

	for index < len(date) {
		if date[index] != '/' {
			currParserVal += string(date[index])
		} else {
			parsedInt, err := strconv.Atoi(currParserVal)
			if err != nil {
				log.Println("Encountered error while parsing the date. Error:	", err)
				return models.Datetime{}, err
			}

			if valToInsert == Date {
				parsedDate.Day = parsedInt
				valToInsert = Month
			} else {
				parsedDate.Month = parsedInt
			}
			currParserVal = ""
		}
		index++
	}

	parsedInt, err := strconv.Atoi(currParserVal)
	if err != nil {
		log.Println("Encountered error while parsing the date. Error:	", err)
		return models.Datetime{}, err
	}
	parsedDate.Year = parsedInt

	if err := FeasibleDate(parsedDate); err != nil {
		return models.Datetime{}, err
	}
	return parsedDate, nil
}

func ParseDateToString(date models.Datetime) string {
	dateInStr := ""
	dateInStr += strconv.Itoa(date.Day) + string('/')
	dateInStr += strconv.Itoa(date.Month) + string('/')
	dateInStr += strconv.Itoa(date.Year)
	return dateInStr
}
