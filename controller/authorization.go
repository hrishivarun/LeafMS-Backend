package controller

import (
	models "LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
	"encoding/json"
	"net/http"
)

// ============================================================================
// ============================================================================
// handle `login`
// ============================================================================
// ============================================================================
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginReq
	if err := service.DecodeJson(r.Body, &loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(loginReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginInfo := service.ValidateCred(loginReq)
	if loginInfo.Username == "" || loginInfo.Status != http.StatusOK {
		http.Error(w, "Error occured in VALIDATION!!", loginInfo.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", loginInfo.Token)

	response, _ := json.MarshalIndent(loginInfo, "", "	")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `admin login`
// ============================================================================
// ============================================================================
func HandleAdminLogin(w http.ResponseWriter, r *http.Request) {
	var adminLoginReq models.LoginReq
	if err := service.DecodeJson(r.Body, &adminLoginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(adminLoginReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if adminLoginReq.Username != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}
	loginInfo := service.ValidateCred(adminLoginReq)
	if adminLoginReq.Username == "" || loginInfo.Status != http.StatusOK {
		http.Error(w, "Error occured", loginInfo.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", loginInfo.Token)

	response, _ := json.MarshalIndent(loginInfo, "", "	")
	w.Write(response)
}
