package service

import (
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"
	"log"
	"os"
	"time"
)

var (
	JWTSecretKey = []byte(os.Getenv("JWT_SECRET"))
	TokenExpiry  = time.Hour * 24
)

const Admin = "admin"

func FilterHolidaysFromLeaveRequest(leaveApplication models.MetaLeaveInfo) (models.MetaLeaveInfo, error) {
	var splitLeaves []models.LeaveInfo
	for _, leave := range leaveApplication.Leaves {
		leaveSlices, err := utils.RemoveHolidayFromLeaveData(leave)
		if err != nil {
			log.Println("Could not remove the holidays from the leave applied. Err : ", err)
			return models.MetaLeaveInfo{}, err
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
