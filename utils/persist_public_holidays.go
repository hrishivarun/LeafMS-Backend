package utils

import (
	db "LeafMS-BackEnd/database"
	models "LeafMS-BackEnd/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var database = db.ConnectDB()

func PersistPublicHolidays(year int, countryCode string) {
	var url = fmt.Sprintf("https://calendarific.com/api/v2/holidays?&api_key=uhXXRzt1AhCbm9h6MKzfqwCU7kT4XFEH&country=%s&year=%d", countryCode, year)
	var client = &http.Client{Timeout: 10 * time.Second}
	var resp, err = client.Get(url)

	if err != nil {
		fmt.Printf("Failed to fetch holidays: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received status code %d\n", resp.StatusCode)
		return
	}

	var holidaysJson models.HolidayApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&holidaysJson); err != nil {
		log.Printf("Error decoding response: %v\n", err)
		return
	}

	var holidaysArr []interface{}
	for _, holidays := range holidaysJson.Response.Holidays {
		holidaysArr = append(holidaysArr, holidays)
	}

	result, err := database.InsertMany("publicHolidays", holidaysArr)
	if err != nil {
		log.Fatalln("Could not persist public holiday data in database!!\n\n Error:=	", err)
		return
	}

	log.Printf("Public holidays successfully inserted!!\n\n %v", result)
}
