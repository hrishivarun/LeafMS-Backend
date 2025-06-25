package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ============================================================================
// ============================================================================
// `leave apply`
// ============================================================================
// ============================================================================
func ApplyForLeave(leaveApplication models.MetaLeaveInfo) (*mongo.UpdateResult, error) {
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
						{Key: "$in", Value: cancelLeaveReq.LeaveIds}}}}}}}})
	if err != nil {
		return nil, err
	}

	return cancellationResult, nil
}

// ============================================================================
// ============================================================================
// `view team's leaves`
// ============================================================================
// ============================================================================
func ViewTeamLeaveInfo(teamName string) ([]models.MetaLeaveInfo, error) {
	teamPeepsRaw, err := database.DbConn.Find("employees", bson.D{
		{Key: "team", Value: teamName}})

	if err != nil {
		return nil, err
	}
	teamPeeps := utils.ReturnEmployees(teamPeepsRaw)
	var peepsUsername []string
	for _, peep := range teamPeeps {
		peepsUsername = append(peepsUsername, peep.Username)
	}

	data, err := database.DbConn.Find("leaves", bson.D{
		{Key: "username", Value: bson.D{{Key: "$in", Value: peepsUsername}}}})

	if err != nil {
		return nil, err
	}

	leaves := utils.ReturnLeaves(data)
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

	holidays := utils.ReturnHolidays(holidaysBson)
	return holidays, nil
}
