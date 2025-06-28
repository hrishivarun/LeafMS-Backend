package service

import (
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	JWTSecretKey = []byte(os.Getenv("JWT_SECRET"))
	TokenExpiry  = time.Hour * 24
)

const Admin = "admin"

var reqValidator = validator.New()

func FilterHolidaysFromLeaveRequest(leaveApplication models.LeaveApplication) (models.LeaveApplication, error) {
	var splitLeaves []models.LeaveInfo
	for _, leave := range leaveApplication.Leaves {
		leaveSlices, err := utils.RemoveHolidayFromLeaveData(leave)
		if err != nil {
			log.Println("Could not remove the holidays from the leave applied. Err : ", err)
			return models.LeaveApplication{}, err
		}

		splitLeaves = append(splitLeaves, leaveSlices...)
	}

	var leavesLackingWeekend []models.LeaveInfo
	for _, leave := range splitLeaves {
		leaveSlices, err := utils.RemoveWeekendsFromLeaveData(leave)
		if err != nil {
			log.Fatalln("There was an error while removing weekends from the applied leave. Err : ", err)
		}

		leavesLackingWeekend = append(leavesLackingWeekend, leaveSlices...)
	}

	leaveApplication.Leaves = leavesLackingWeekend
	return leaveApplication, nil
}

func CreateViewLeavesFilter(viewLeavesReq models.ViewLeavesReq) bson.D {
	filter := bson.D{{Key: "username", Value: viewLeavesReq.Username}}

	if viewLeavesReq.Year != 0 && viewLeavesReq.Month != 0 {
		filter = append(filter, bson.E{Key: "leaves", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
			{Key: "startDate", Value: bson.D{
				{Key: "$gte", Value: time.Date(viewLeavesReq.Year, time.Month(viewLeavesReq.Month), 1, 0, 0, 0, 0, time.Local)},
				{Key: "$lt", Value: time.Date(viewLeavesReq.Year, time.Month(viewLeavesReq.Month)+1, 1, 0, 0, 0, 0, time.Local)},
			}},
		}}}})
	} else if viewLeavesReq.Year != 0 {
		filter = append(filter, bson.E{Key: "leaves", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
			{Key: "startDate", Value: bson.D{
				{Key: "$gte", Value: time.Date(viewLeavesReq.Year, 1, 1, 0, 0, 0, 0, time.Local)},
				{Key: "$lt", Value: time.Date(viewLeavesReq.Year+1, 1, 1, 0, 0, 0, 0, time.Local)},
			}},
		}}}})
	}
	return filter
}

func DecodeJson(body io.ReadCloser, model interface{}) error {
	err := json.NewDecoder(body).Decode(model)
	if err != nil {
		log.Fatal("mar gayo re, sala json decoder phar rha hai")
		return err
	}
	return nil
}

func ValidateRequest(request any) error {
	if err := reqValidator.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errorMessages := ""
		for _, fieldErr := range validationErrors {
			errorMessages += fmt.Sprintf(fieldErr.Field() + "is required\n")
		}
		return errors.New(errorMessages)
	}
	return nil
}
