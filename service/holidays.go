package service

import (
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================
// ============================================================================
// `view holidays`
// ============================================================================
// ============================================================================
func ViewHolidays(filter models.HolidaysFilter) ([]models.Holiday, error) {
	// Build query dynamically
	query := bson.D{
		{Key: "country.id", Value: filter.Country},
		{Key: "date.datetime.year", Value: filter.Year},
	}

	// Only add month if provided (non-zero)
	if filter.Month != 0 {
		query = append(query, bson.E{Key: "date.datetime.month", Value: filter.Month})
	}

	holidaysBson, err := dbConn.Find("publicHolidays", query)
	if err != nil {
		return nil, err
	}

	holidays := utils.ReturnHolidays(holidaysBson)
	return holidays, nil
}
