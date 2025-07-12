package service

import (
	"LeafMS-BackEnd/database"
	models "LeafMS-BackEnd/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PersistPublicHolidays(year int, countryCode string) (*mongo.InsertManyResult, error) {
	var url = fmt.Sprintf("https://calendarific.com/api/v2/holidays?&api_key=uhXXRzt1AhCbm9h6MKzfqwCU7kT4XFEH&country=%s&year=%d", countryCode, year)
	var client = &http.Client{Timeout: 10 * time.Second}
	var resp, err = client.Get(url)

	if err != nil {
		fmt.Printf("Failed to fetch holidays: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		return nil, err
	}

	var holidaysJson models.HolidayApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&holidaysJson); err != nil {
		log.Printf("Error decoding response: %v\n", err)
		return nil, err
	}

	var holidaysArr []any
	for _, holidays := range holidaysJson.Response.Holidays {
		holidaysArr = append(holidaysArr, holidays)
	}

	result, err := database.DbConn.InsertMany("publicHolidays", holidaysArr)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("Could not persist public holiday data in database!!\n\n Error:=	", err)
		return nil, err
	}

	log.Printf("Public holidays successfully inserted!!\n\n %v", result)
	return result, nil
}

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

	holidaysBson, err := database.DbConn.Find("publicHolidays", query)
	if err != nil {
		return nil, err
	}

	holidays := database.ConvertRawBsonToHolidays(holidaysBson)
	return holidays, nil
}
