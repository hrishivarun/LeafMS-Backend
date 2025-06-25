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
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var filter models.ViewApplications
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if username != filter.ApproverName {
		w.WriteHeader(http.StatusUnauthorized)
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
	username, ok := r.Context().Value("username").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var leaveApplications models.ApproveLeaveReq
	if err := json.NewDecoder(r.Body).Decode(&leaveApplications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if username != leaveApplications.Username {
		w.WriteHeader(http.StatusUnauthorized)
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
