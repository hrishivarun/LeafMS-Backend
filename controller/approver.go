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
	var filter models.ViewApplications
	if err := service.JsonDecoderWrapper(r.Body, &filter); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != filter.ApproverName {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusUnauthorized)
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
	var leaveApplications models.ApproveLeaveReq
	if err := service.JsonDecoderWrapper(r.Body, &leaveApplications); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); !service.ValidateApprover(leaveApplications.Username, loggedUsername) {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusUnauthorized)
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
