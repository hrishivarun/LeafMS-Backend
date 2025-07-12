package controller

import (
	models "LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
	"encoding/json"
	"net/http"
)

// ============================================================================
// ============================================================================
// handle `view leave applications`
// ============================================================================
// ============================================================================
func HandleViewLeaveApplications(w http.ResponseWriter, r *http.Request) {
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
	if !service.ValidateApprover(viewReq.Username, loggedUsername) {
		http.Error(w, "You do not have access to check this user's information! stop messing bruv", http.StatusUnauthorized)
		return
	}

	data, err := service.ViewLeaveApplications(viewReq, loggedUsername)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, _ := json.MarshalIndent(data, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `leaves approval`
// ============================================================================
// ============================================================================
func HandleLeaveApproval(w http.ResponseWriter, r *http.Request) {
	var leaveApplications models.ResolveLeaveReq
	if err := service.DecodeJson(r.Body, &leaveApplications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(leaveApplications); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); !service.ValidateApprover(leaveApplications.Username, loggedUsername) {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusUnauthorized)
		return
	}

	result, err := service.ResolveLeave(leaveApplications)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, _ := json.MarshalIndent(result, "", "	")
	w.Write(response)
}
