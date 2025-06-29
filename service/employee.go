package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============================================================================
// ============================================================================
// `leave apply`
// ============================================================================
// ============================================================================
func ApplyForLeave(leaveApplication models.LeaveApplication) (*mongo.UpdateResult, error) {
	filteredLeaves, err := FilterHolidaysFromLeaveRequest(leaveApplication)
	if err != nil {
		return nil, err
	}

	result, err := database.DbConn.UpdateOne("leaves",
		bson.D{
			{Key: "username", Value: filteredLeaves.Username}},
		bson.D{
			{Key: "$push", Value: bson.D{
				{Key: "leaves", Value: bson.D{
					{Key: "$each", Value: filteredLeaves.Leaves}}}}}})

	if err != nil {
		log.Println("Encountered error while persisting applied leaves in Database. Err : ", err)
		return nil, err
	}
	return result, nil
}

// ============================================================================
// ============================================================================
// `cancel leaves`
// ============================================================================
// ============================================================================
func CancelLeave(cancelLeaveReq models.CancelLeavesReq) (*mongo.UpdateResult, error) {
	cancellationResult, err := database.DbConn.UpdateOne("leaves",
		bson.D{
			{Key: "username", Value: cancelLeaveReq.Username}},
		bson.D{
			{Key: "$pull", Value: bson.D{
				{Key: "leaves", Value: bson.D{
					{Key: "id", Value: bson.D{
						{Key: "$in", Value: cancelLeaveReq.LeaveIds}}},
					{Key: "status", Value: bson.D{
						{Key: "$nin", Value: bson.A{models.Rejected, models.Cancelled}}}}}}}}})
	if err != nil {
		return nil, err
	}

	return cancellationResult, nil
}

// ============================================================================
// ============================================================================
// `view leaves`
// ============================================================================
// ============================================================================
func ViewLeaves(viewReq models.ViewLeavesReq) ([]models.LeaveInfo, error) {
	viewFilter := CreateViewLeavesFilter(viewReq)
	data, err := database.DbConn.Find("leaves", viewFilter)
	if err != nil {
		return nil, err
	}
	if data == nil {
		log.Fatal("No Leave entry found for given query")
		return nil, nil
	}

	return database.ConvertRawBsonToLeaves(data), nil
}

// ============================================================================
// ============================================================================
// `view team's leaves`
// ============================================================================
// ============================================================================
func ViewTeamLeaveInfo(employeeUsername string) ([]models.LeaveInfo, error) {
	userInfoRaw, err := database.DbConn.FindOne("employees", bson.D{{Key: "username", Value: employeeUsername}})
	if err != nil {
		return nil, err
	}
	var userInfo models.Employee
	if err := bson.Unmarshal(userInfoRaw, &userInfo); err != nil {
		log.Fatal("The decoding of employee from raw bson document failed!\nError:-\n\n", err)
		return nil, err
	}
	teamPeepsRaw, err := database.DbConn.Find("employees", bson.D{
		{Key: "team", Value: userInfo.Team}})

	if err != nil {
		return nil, err
	}
	teamPeeps := database.ConvertRawBsonToEmployees(teamPeepsRaw)
	var peepsUsername []string
	for _, peep := range teamPeeps {
		peepsUsername = append(peepsUsername, peep.Username)
	}

	data, err := database.DbConn.Find("leaves", bson.D{
		{Key: "username", Value: bson.D{{Key: "$in", Value: peepsUsername}}}})

	if err != nil {
		return nil, err
	}

	leaves := database.ConvertRawBsonToLeaves(data)
	return leaves, nil
}

// ============================================================================
// ============================================================================
// `view holidays`
// ============================================================================
// ============================================================================
func ViewHolidays(filter models.HolidaysFilter) ([]models.Holiday, error) {
	// Build query dynamically
	query := bson.D{
		{Key: "country.id", Value: filter.Country},
		{Key: "date.datetime.year", Value: filter.Year},
	}

	// Only add month if provided (non-zero)
	if filter.Month != 0 {
		query = append(query, bson.E{Key: "date.datetime.month", Value: filter.Month})
	}

	holidaysBson, err := database.DbConn.Find("publicHolidays", query)
	if err != nil {
		return nil, err
	}

	holidays := database.ConvertRawBsonToHolidays(holidaysBson)
	return holidays, nil
}
