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
	var user models.Employee
	if err := service.DecodeJson(r.Body, &user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authenticatedUser, loginInfo := service.ValidateCred(user)
	if authenticatedUser.Username == "" || loginInfo.Status != http.StatusOK {
		http.Error(w, "Error occured", loginInfo.Status)
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
	var user models.Employee
	if err := service.DecodeJson(r.Body, &user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if user.Username != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}
	authenticatedUser, loginInfo := service.ValidateCred(user)
	if authenticatedUser.Username == "" || loginInfo.Status != http.StatusOK {
		http.Error(w, "Error occured", loginInfo.Status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", loginInfo.Token)

	response, _ := json.MarshalIndent(loginInfo, "", "	")
	w.Write(response)
}
