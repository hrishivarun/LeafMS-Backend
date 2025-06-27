package controller

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
	"LeafMS-BackEnd/utils"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

// ============================================================================
// ============================================================================
// handle `view employees`
// ============================================================================
// ============================================================================
func HandleViewEmployees(w http.ResponseWriter, r *http.Request) {
	var admin models.Employee
	if err := service.JsonDecoderWrapper(r.Body, &admin); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	records, err := database.DbConn.Find("employees", bson.D{})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	response, _ := json.MarshalIndent(records, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `register employee`
// ============================================================================
// ============================================================================
func HandleRegisterEmployee(w http.ResponseWriter, r *http.Request) {
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	var newEmployee models.Employee
	if err := service.JsonDecoderWrapper(r.Body, &newEmployee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	registrationResult, err := database.DbConn.InsertOne("employees", newEmployee)
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if registrationResult.InsertedID == 0 {
		http.Error(w, "For some reason, the employee creation did not go through?", http.StatusInternalServerError)
	}
	registrationResponse, _ := json.MarshalIndent(registrationResult, "", " ")
	w.Write(registrationResponse)
}

// ============================================================================
// ============================================================================
// handle `remove employee`
// ============================================================================
// ============================================================================
func HandleRemoveEmployee(w http.ResponseWriter, r *http.Request) {
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	var employeeToDelete models.Employee
	if err := service.JsonDecoderWrapper(r.Body, &employeeToDelete); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletionResult, err := database.DbConn.DeleteOne("employees", bson.D{{Key: "username", Value: employeeToDelete.Username}})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if deletionResult.DeletedCount == 0 {
		http.Error(w, "Wtf bro? no employee was deleted, what are you doin?? you sure this is the username", http.StatusNotFound)
	}
	deletionResponse, _ := json.MarshalIndent(deletionResult, "", " ")
	w.Write(deletionResponse)
}

// ============================================================================
// ============================================================================
// handle `post holidays`
// ============================================================================
// ============================================================================
func HandlePostHolidays(w http.ResponseWriter, r *http.Request) {
	var postHolidayReq models.PostHoliday
	if err := service.JsonDecoderWrapper(r.Body, &postHolidayReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	postHolidayRes, err := utils.PersistPublicHolidays(postHolidayReq.Year, postHolidayReq.Country)
	if err != nil {
		http.Error(w, "oopsie, holiday write request did not GO THORUGH SIR!!!", http.StatusInternalServerError)
		return
	}
	jsonRes := models.InsertManyResult{InsertedIDs: postHolidayRes.InsertedIDs}
	res, _ := json.MarshalIndent(jsonRes, "", " ")
	w.Write(res)
}
