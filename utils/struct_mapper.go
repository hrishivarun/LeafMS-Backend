package utils

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"

	models "LeafMS-BackEnd/models"
)

func ReturnLeaves(data []bson.Raw) []models.MetaLeaveInfo {
	var leaves []models.MetaLeaveInfo
	for _, entry := range data {
		var leave models.MetaLeaveInfo
		if err := bson.Unmarshal(entry, &leave); err != nil {
			log.Fatal(
				"The decoding of leaveApplication from raw bson document failed!\nError:-\n\n", err)
		}
		leaves = append(leaves, leave)
	}
	return leaves
}

func ReturnEmployees(data []bson.Raw) []models.Employee {
	var employees []models.Employee
	for _, entry := range data {
		var employee models.Employee
		if err := bson.Unmarshal(entry, &employee); err != nil {
			log.Fatal(
				"The decoding of employee from raw bson document failed!\nError:-\n\n", err)
		}
		employees = append(employees, employee)
	}
	return employees
}

func ReturnHolidays(data []bson.Raw) []models.Holiday {
	var holidays []models.Holiday
	for _, entry := range data {
		var holiday models.Holiday
		if err := bson.Unmarshal(entry, &holiday); err != nil {
			log.Fatal(
				"The decoding of employee from raw bson document failed!\nError:-\n\n", err)
		}
		holidays = append(holidays, holiday)
	}
	return holidays
}
