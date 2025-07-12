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
	"go.mongodb.org/mongo-driver/mongo"
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
			log.SetPrefix("WARNING: ")
			log.Println("There was an error while removing weekends from the applied leave. Err : ", err)
		}

		leavesLackingWeekend = append(leavesLackingWeekend, leaveSlices...)
	}

	leaveApplication.Leaves = leavesLackingWeekend
	return leaveApplication, nil
}

// func CreateViewLeavesFilter(viewLeavesReq models.ViewApplicationsReq) bson.D {
// 	filter := bson.D{{Key: "username", Value: viewLeavesReq.Username}}

//		if viewLeavesReq.Year != nil && viewLeavesReq.Month != nil {
//			filter = append(filter, bson.E{Key: "leaves", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
//				{Key: "startDate", Value: bson.D{
//					{Key: "$gte", Value: time.Date(*viewLeavesReq.Year, time.Month(*viewLeavesReq.Month), 1, 0, 0, 0, 0, time.Local)},
//					{Key: "$lt", Value: time.Date(*viewLeavesReq.Year, time.Month(*viewLeavesReq.Month)+1, 1, 0, 0, 0, 0, time.Local)},
//				}},
//			}}}})
//		} else if *viewLeavesReq.Year != 0 {
//			filter = append(filter, bson.E{Key: "leaves", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
//				{Key: "startDate", Value: bson.D{
//					{Key: "$gte", Value: time.Date(*viewLeavesReq.Year, 1, 1, 0, 0, 0, 0, time.Local)},
//					{Key: "$lt", Value: time.Date(*viewLeavesReq.Year+1, 1, 1, 0, 0, 0, 0, time.Local)},
//				}},
//			}}}})
//		}
//		return filter
//	}
func AddUsernameFilterToDbPipeline(pipeline mongo.Pipeline, usernames []string) mongo.Pipeline {
	pipeline = append(pipeline,
		bson.D{{Key: "$match", Value: bson.D{{Key: "username", Value: bson.D{{Key: "$in", Value: usernames}}}}}})
	return pipeline
}

func AddApproverFilterToDbPipeline(pipeline mongo.Pipeline, approverUsername string) mongo.Pipeline {
	pipeline = append(pipeline,
		bson.D{{Key: "$match", Value: bson.D{{Key: "approver", Value: approverUsername}}}})
	return pipeline
}

func ComposeLeaveFilter(req models.ViewApplicationsReq) bson.D {
	// 2. Build conditions for $filter
	var leavesConds bson.A

	if req.LeaveType != nil {
		leavesConds = append(leavesConds,
			bson.D{{Key: "$eq", Value: bson.A{"$$leave.type", *req.LeaveType}}},
		)
	}
	if req.Status != nil {
		leavesConds = append(leavesConds,
			bson.D{{Key: "$eq", Value: bson.A{"$$leave.status", *req.Status}}},
		)
	}
	if req.Year != nil {
		leavesConds = append(leavesConds,
			bson.D{{Key: "$eq", Value: bson.A{
				bson.D{{Key: "$year", Value: "$$leave.startDate"}},
				*req.Year,
			}}},
		)
	}
	if req.Month != nil {
		leavesConds = append(leavesConds,
			bson.D{{Key: "$eq", Value: bson.A{
				bson.D{{Key: "$month", Value: "$$leave.startDate"}},
				*req.Month,
			}}},
		)
	}

	// Compose the filter condition
	var composedCond bson.D
	if len(leavesConds) == 1 {
		composedCond = leavesConds[0].(bson.D)
	} else if len(leavesConds) > 1 {
		composedCond = bson.D{{Key: "$and", Value: leavesConds}}
	} else {
		composedCond = bson.D{} // Always true, so no filter
	}
	return composedCond
}

func BuildLeaveFilterPipeline(pipeline mongo.Pipeline, req models.ViewApplicationsReq) mongo.Pipeline {
	cond := ComposeLeaveFilter(req)

	// Project filtered leaves array
	pipeline = append(pipeline,
		bson.D{{
			Key: "$project", Value: bson.D{
				{Key: "leaves", Value: bson.D{
					{Key: "$filter", Value: bson.D{
						{Key: "input", Value: "$leaves"},
						{Key: "as", Value: "leave"},
						{Key: "cond", Value: cond},
					}},
				}},
			},
		}},
	)
	return pipeline
}

func FlattenDbResult(pipeline mongo.Pipeline, newRoot string) mongo.Pipeline {
	// 4. Unwind and flatten
	pipeline = append(pipeline,
		bson.D{{Key: "$unwind", Value: "$leaves"}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: newRoot}}}})
	return pipeline
}

func DecodeJson(body io.ReadCloser, model any) error {
	err := json.NewDecoder(body).Decode(model)
	if err != nil {
		log.SetPrefix("WARNING: ")
		log.Println("mar gayo re, sala json decoder phar rha hai")
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
