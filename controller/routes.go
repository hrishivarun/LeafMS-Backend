package controller

import (
	"encoding/json"
	"log"
	"net/http"

	db "LeafMS-BackEnd/database"
	"LeafMS-BackEnd/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userInfo db.User

// ============================================================================
// ============================================================================
// handle `login`
// ============================================================================
// ============================================================================
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user db.User

	log.Println("started login api")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request payload!!!"))
		return
	}

	//Authenticate the user credentials with the database
	user, loginInfo := validateCred(user)
	userInfo = user //saving userInfo
	log.Println("validated cred")

	sessiondId := uuid.New().String()
	jwtToken, err := generateJWT(sessiondId)
	if err != nil {
		log.Fatalf("couldn't generate JWT auth token.\nError: %v\n", err)
	}
	w.Header().Add("Authorization", jwtToken)
	w.Header().Add("Session-Id", sessiondId)

	response, _ := json.MarshalIndent(loginInfo, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `apply leaves`
// ============================================================================
// ============================================================================
func HandleApply(w http.ResponseWriter, r *http.Request) {
	var leaveApplication db.Leaves
	err := json.NewDecoder(r.Body).Decode(&leaveApplication)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var splitLeaves []db.LeaveData
	for _, leave := range leaveApplication.Leaves {
		leaveSlices, err := utils.RemoveHolidayFromLeaveData(leave)
		if err != nil {
			log.Println("Could not remove the holidays from the leave applied. Err : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		splitLeaves = append(splitLeaves, leaveSlices...)
	}

	var leavesLackingWeekend []db.LeaveData
	for _, leave := range splitLeaves {
		leaveSlices, err := utils.RemoveWeekendsFromLeaveData(leave)
		if err != nil {
			log.Fatalln("There was an error while removing weekends from the applied leave. Err : ", err)
		}

		leavesLackingWeekend = append(leavesLackingWeekend, leaveSlices...)
	}

	leaveApplication.Leaves = leavesLackingWeekend

	result, err := database.UpdateOne("leaves", bson.D{
		{Key: "username", Value: leaveApplication.Username},
	}, bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "leaves", Value: bson.D{
				{Key: "$each", Value: leaveApplication.Leaves},
			}},
		}},
	})
	if err != nil {
		log.Println("Encountered error while persisting applied leaves in Database. Err : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal("No User with the username: " + leaveApplication.Username + " exists.")
		w.Write(response)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response, _ := json.MarshalIndent(result, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view leaves`
// ============================================================================
// ============================================================================
func HandleViewLeaves(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := database.Find("leaves", bson.D{
		{Key: "username", Value: user.Username}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leaves := utils.ReturnLeaves(data)
	response, _ := json.MarshalIndent(leaves, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view team's leaves`
// ============================================================================
// ============================================================================
func HandleViewTeamLeaves(w http.ResponseWriter, r *http.Request) {
	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == db.User{} || userInfo.Username != user.Username) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}

	teamPeepsRaw, err := database.Find("employees", bson.D{
		{Key: "team", Value: user.Team}})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	teamPeeps := utils.ReturnUsers(teamPeepsRaw)
	var peepsUsername []string
	for _, peep := range teamPeeps {
		peepsUsername = append(peepsUsername, peep.Username)
	}

	data, err := database.Find("leaves", bson.D{
		{Key: "username", Value: bson.D{{Key: "$in", Value: peepsUsername}}}})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leaves := utils.ReturnLeaves(data)
	response, _ := json.MarshalIndent(leaves, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view leave applications`
// ============================================================================
// ============================================================================
func HandleViewLeaveApplications(w http.ResponseWriter, r *http.Request) {
	var filter ViewApplications
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var pipeline mongo.Pipeline
	// Always filter by approver
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
		{Key: "approver", Value: filter.ApproverName},
	}}})

	// If IsLeaveAprroved is provided, filter by it
	if filter.IsLeaveAprroved != nil {
		pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "leaves", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$leaves"},
					{Key: "as", Value: "leave"},
					{Key: "cond", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$$leave.approved", *filter.IsLeaveAprroved}},
					}},
				}},
			}},
		}}})
	}

	data, err := database.Aggregate("leaves", pipeline)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	leaveApplications := utils.ReturnLeaves(data)
	response, _ := json.MarshalIndent(leaveApplications, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `leaves approval`
// ============================================================================
// ============================================================================
func HandleLeaveApproval(w http.ResponseWriter, r *http.Request) {
	var leaveData db.Leaves
	if err := json.NewDecoder(r.Body).Decode(&leaveData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedResult, err := database.UpdateOne("leaves", bson.D{
		{Key: "username", Value: leaveData.Username}, {
			Key: "leaves", Value: bson.D{{
				Key: "$elemMatch", Value: bson.D{{Key: "id", Value: leaveData.Leaves[0].Id}}}}}, //possible bug, why only matching for Leaves[0], why not for other IDs?
	}, bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "leaves.$.approved", Value: leaveData.Leaves[0].Approved},
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(updatedResult, "", "	")
	w.Write(response)

}

func HandleViewHolidays(w http.ResponseWriter, r *http.Request) {
	var holidayArgs db.HolidayArgs
	if err := json.NewDecoder(r.Body).Decode(&holidayArgs); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	holidaysBson, err := database.Find("publicHolidays", bson.D{
		{Key: "country.id", Value: holidayArgs.Country},
		{Key: "date.datetime.year", Value: holidayArgs.Year},
	})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	holidays := utils.ReturnHolidays(holidaysBson)

	serverRes, _ := json.MarshalIndent(holidays, "", "	")
	w.Write(serverRes)
}
