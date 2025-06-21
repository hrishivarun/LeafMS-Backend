package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"LeafMS-BackEnd/database"
	models "LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
	"LeafMS-BackEnd/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var userInfo models.Employee

// ============================================================================
// ============================================================================
// handle `login`
// ============================================================================
// ============================================================================
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var user models.Employee

	log.Println("started login api")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request payload!!!"))
		return
	}

	//Authenticate the user credentials with the database
	user, loginInfo := service.ValidateCred(user)
	userInfo = user //saving userInfo
	log.Println("validated cred")

	sessiondId := uuid.New().String()
	jwtToken, err := service.GenerateJWT(sessiondId)
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
// handle `leave apply`
// ============================================================================
// ============================================================================
func HandleApply(w http.ResponseWriter, r *http.Request) {
	var leaveApplication models.MetaLeaveInfo
	err := json.NewDecoder(r.Body).Decode(&leaveApplication)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == models.Employee{} || userInfo.Username != leaveApplication.Username) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	result, err := service.ApplyForLeave(leaveApplication)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal("No Employee with the username: " + leaveApplication.Username + " exists.")
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
	var user models.Employee
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == models.Employee{} || userInfo.Username != user.Username) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}

	data, err := database.DbConn.Find("leaves", bson.D{
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
// handle `cancel leaves`
// ============================================================================
// ============================================================================
func HandleCancelLeave(w http.ResponseWriter, r *http.Request) {
	var cancelLeaveReq models.CancelLeavesReq

	err := json.NewDecoder(r.Body).Decode(&cancelLeaveReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == models.Employee{} || userInfo.Username != cancelLeaveReq.Username) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}

	cancellationRes, err := service.CancelLeave(cancelLeaveReq)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if cancellationRes.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal("No leave application with given leave ID")
		w.Write(response)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	response, _ := json.MarshalIndent(cancellationRes, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view team's leaves`
// ============================================================================
// ============================================================================
func HandleViewTeamLeaves(w http.ResponseWriter, r *http.Request) {
	var user models.Employee
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == models.Employee{} || userInfo.Username != user.Username) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}

	leaves, err := service.ViewTeamLeaveInfo(user.Team)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if len(leaves) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, _ := json.MarshalIndent(leaves, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view leave applications`
// ============================================================================
// ============================================================================
func HandleViewLeaveApplications(w http.ResponseWriter, r *http.Request) {
	var filter models.ViewApplications
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (userInfo == models.Employee{} || userInfo.Username != filter.ApproverName) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}
	applications, err := service.ViewLeaveApplications(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, _ := json.MarshalIndent(applications, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `leaves approval`
// ============================================================================
// ============================================================================
func HandleLeaveApproval(w http.ResponseWriter, r *http.Request) {
	var leaveApplications models.MetaLeaveInfo
	if err := json.NewDecoder(r.Body).Decode(&leaveApplications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if (userInfo == models.Employee{} || userInfo.Username != leaveApplications.Approver) {
		w.WriteHeader((http.StatusUnauthorized))
		return
	}
	updatedResult, err := service.ApproveLeave(leaveApplications)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(updatedResult, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view holidays`
// ============================================================================
// ============================================================================
func HandleViewHolidays(w http.ResponseWriter, r *http.Request) {
	var filter models.HolidaysFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	holidays, err := service.ViewHolidays(filter)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	serverRes, _ := json.MarshalIndent(holidays, "", "	")
	w.Write(serverRes)
}
