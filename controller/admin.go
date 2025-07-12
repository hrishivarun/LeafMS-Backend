package controller

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/service"
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
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	var employeesInfoReq models.EmployeesCrudReq
	if err := service.DecodeJson(r.Body, &employeesInfoReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(employeesInfoReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rawResult, err := database.DbConn.Find("employees", bson.D{{Key: "username", Value: bson.D{{Key: "$in", Value: employeesInfoReq.Usernames}}}})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	result := database.ConvertRawBsonToEmployees(rawResult)
	response, _ := json.MarshalIndent(result, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `view all employees`
// ============================================================================
// ============================================================================
func HandleViewAllEmployees(w http.ResponseWriter, r *http.Request) {
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	rawResult, err := database.DbConn.Find("employees", bson.D{})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	result := database.ConvertRawBsonToEmployees(rawResult)
	response, _ := json.MarshalIndent(result, "", " ")
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
	if err := service.DecodeJson(r.Body, &newEmployee); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(newEmployee); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DbConn.InsertOne("employees", newEmployee)
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if result.InsertedID == 0 {
		http.Error(w, "For some reason, the employee creation did not go through?", http.StatusInternalServerError)
	}

	leaveDocCreationResult, err := database.DbConn.InsertOne("leaves", models.LeaveDoc{Username: newEmployee.Username, Leaves: []models.LeaveInfo{}})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if leaveDocCreationResult.InsertedID == 0 {
		http.Error(w, "For some reason, the employee creation did not go through?", http.StatusInternalServerError)
	}
	response, _ := json.MarshalIndent(result, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `remove employees`
// ============================================================================
// ============================================================================
func HandleRemoveEmployees(w http.ResponseWriter, r *http.Request) {
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}

	var employeesToDelete models.EmployeesCrudReq
	if err := service.DecodeJson(r.Body, &employeesToDelete); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(employeesToDelete); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := database.DbConn.DeleteMany("employees", bson.D{{Key: "username", Value: bson.D{{Key: "$in", Value: employeesToDelete.Usernames}}}})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Wtf bro? no employee was deleted, what are you doin?? you sure this is the username", http.StatusNotFound)
	}

	leaveDocDeletionRes, err := database.DbConn.DeleteMany("leaves", bson.D{{Key: "username", Value: bson.D{{Key: "$in", Value: employeesToDelete.Usernames}}}})
	if err != nil {
		http.Error(w, "Something shitty happened here!!!", http.StatusInternalServerError)
	}
	if leaveDocDeletionRes.DeletedCount == 0 {
		http.Error(w, "Wtf bro? no employee was deleted, what are you doin?? you sure this is the username", http.StatusNotFound)
	}
	response, _ := json.MarshalIndent(result, "", " ")
	w.Write(response)
}

// ============================================================================
// ============================================================================
// handle `post holidays`
// ============================================================================
// ============================================================================
func HandlePostHolidays(w http.ResponseWriter, r *http.Request) {
	if loggedUsername, _ := r.Context().Value("username").(string); loggedUsername != service.Admin {
		http.Error(w, "You do not have access to this API! stop messing bruv", http.StatusForbidden)
		return
	}
	var postHolidayReq models.PostHoliday
	if err := service.DecodeJson(r.Body, &postHolidayReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := service.ValidateRequest(postHolidayReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	postHolidayRes, err := service.PersistPublicHolidays(postHolidayReq.Year, postHolidayReq.Country)
	if err != nil {
		http.Error(w, "oopsie, holiday write request did not GO THORUGH SIR!!!", http.StatusInternalServerError)
		return
	}
	result := models.InsertManyResult{InsertedIDs: postHolidayRes.InsertedIDs}
	res, _ := json.MarshalIndent(result, "", " ")
	w.Write(res)
}
