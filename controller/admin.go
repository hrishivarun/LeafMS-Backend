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
