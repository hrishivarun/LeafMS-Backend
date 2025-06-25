package controller

import (
	"encoding/json"
	"net/http"

	"LeafMS-BackEnd/database"
	models "LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
	"LeafMS-BackEnd/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================
// ============================================================================
// handle `leave apply`
// ============================================================================
// ============================================================================
func HandleApply(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var leaveApplication models.MetaLeaveInfo
	err := json.NewDecoder(r.Body).Decode(&leaveApplication)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if username != leaveApplication.Username {
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
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user models.Employee
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// relic of a past better not talked about

	if username != user.Username {
		w.WriteHeader(http.StatusUnauthorized)
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
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var cancelLeaveReq models.CancelLeavesReq

	err := json.NewDecoder(r.Body).Decode(&cancelLeaveReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if username != cancelLeaveReq.Username {
		w.WriteHeader(http.StatusUnauthorized)
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
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var user models.Employee
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if username != user.Username {
		w.WriteHeader(http.StatusUnauthorized)
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
