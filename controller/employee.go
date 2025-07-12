package controller

import (
	"encoding/json"
	"net/http"

	models "LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
)

// ============================================================================
// ============================================================================
// handle `leave apply`
// ============================================================================
// ============================================================================
func HandleApply(w http.ResponseWriter, r *http.Request) {
	var leaveApplication models.LeaveApplication
	if err := service.DecodeJson(r.Body, &leaveApplication); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(leaveApplication); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != leaveApplication.Username {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
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
// handle `cancel leaves`
// ============================================================================
// ============================================================================
func HandleCancelLeave(w http.ResponseWriter, r *http.Request) {
	var cancelLeaveReq models.CancelLeavesReq
	if err := service.DecodeJson(r.Body, &cancelLeaveReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(cancelLeaveReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != cancelLeaveReq.Username {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	result, err := service.CancelLeave(cancelLeaveReq)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal("No leave application with given leave ID exists for the user.")
		w.Write(response)
		return
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
	var viewReq models.ViewApplicationsReq
	if err := service.DecodeJson(r.Body, &viewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(viewReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != *viewReq.Username {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	data, err := service.ViewLeaves(viewReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, _ := json.MarshalIndent(data, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view team's leaves`
// ============================================================================
// ============================================================================
func HandleViewTeamLeaves(w http.ResponseWriter, r *http.Request) {
	var viewReq models.ViewApplicationsReq
	if err := service.DecodeJson(r.Body, &viewReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(viewReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loggedUsername, _ := r.Context().Value("username").(string)

	data, err := service.ViewTeamLeaveInfo(loggedUsername, viewReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if len(data) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, _ := json.MarshalIndent(data, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view holidays`
// ============================================================================
// ============================================================================
func HandleViewHolidays(w http.ResponseWriter, r *http.Request) {
	var filter models.HolidaysFilter
	if err := service.DecodeJson(r.Body, &filter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(filter); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
