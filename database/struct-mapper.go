package database

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"

	models "LeafMS-BackEnd/models"
)

func ConvertRawBsonToLeaves(data []bson.Raw) []models.LeaveInfo {
	var leaves []models.LeaveInfo
	for _, entry := range data {
		var leave models.LeaveInfo
		if err := bson.Unmarshal(entry, &leave); err != nil {
			log.SetPrefix("WARNING: ")
			log.Println(
				"The decoding of leaveApplication from raw bson document failed!\nError:-\n\n", err)
		}
		leaves = append(leaves, leave)
	}
	return leaves
}

func ConvertRawBsonToEmployees(data []bson.Raw) []models.Employee {
	var employees []models.Employee
	for _, entry := range data {
		var employee models.Employee
		if err := bson.Unmarshal(entry, &employee); err != nil {
			log.SetPrefix("WARNING: ")
			log.Println("The decoding of employee from raw bson document failed!\nError:-\n\n", err)
		}
		employees = append(employees, employee)
	}
	return employees
}

func ConvertRawBsonToHolidays(data []bson.Raw) []models.Holiday {
	var holidays []models.Holiday
	for _, entry := range data {
		var holiday models.Holiday
		if err := bson.Unmarshal(entry, &holiday); err != nil {
			log.SetPrefix("WARNING: ")
			log.Println(
				"The decoding of employee from raw bson document failed!\nError:-\n\n", err)
		}
		holidays = append(holidays, holiday)
	}
	return holidays
}
