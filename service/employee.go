package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	if len(filteredLeaves.Leaves) > 0 {
		opts := options.Update().SetUpsert(true)
		result, err := database.DbConn.UpdateOne("leaves",
			bson.D{
				{Key: "username", Value: filteredLeaves.Username}},
			bson.D{
				{Key: "$push", Value: bson.D{
					{Key: "leaves", Value: bson.D{
						{Key: "$each", Value: filteredLeaves.Leaves}}}}}}, opts)

		if err != nil {
			log.Println("Encountered error while persisting applied leaves in Database. Err : ", err)
			return nil, err
		}
		return result, nil
	}
	return nil, nil
}

// ============================================================================
// ============================================================================
// `cancel leaves`
// ============================================================================
// ============================================================================
func CancelLeave(cancelLeaveReq models.CancelLeavesReq) (*mongo.UpdateResult, error) {
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "leaves.$[elem].status", Value: models.Cancelled},
		}},
	}
	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []any{
			bson.D{{Key: "elem.id", Value: bson.D{{Key: "$in", Value: cancelLeaveReq.LeaveIds}}}},
			bson.D{{Key: "elem.status", Value: bson.D{{Key: "$nin", Value: bson.A{models.Rejected, models.Cancelled}}}}},
		},
	})

	cancellationResult, err := database.DbConn.UpdateOne("leaves", bson.D{{Key: "username", Value: cancelLeaveReq.Username}}, update, arrayFilters)
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
func ViewLeaves(viewReq models.ViewApplicationsReq) ([]models.LeaveInfo, error) {
	var leavePipeline mongo.Pipeline
	if viewReq.Username != nil {
		leavePipeline = AddUsernameFilterToDbPipeline(leavePipeline, []string{*viewReq.Username})
	}
	leavePipeline = BuildLeaveFilterPipeline(leavePipeline, viewReq)
	leavePipeline = FlattenDbResult(leavePipeline, "$leaves")
	data, err := database.DbConn.Aggregate("leaves", leavePipeline)
	if err != nil {
		return nil, err
	}
	if data == nil {
		log.SetPrefix("WARNING: ")
		log.Println("No Leave entry found for given query")
		return nil, nil
	}

	return database.ConvertRawBsonToLeaves(data), nil
}

// ============================================================================
// ============================================================================
// `view team's leaves`
// ============================================================================
// ============================================================================
func ViewTeamLeaveInfo(employeeUsername string, viewReq models.ViewApplicationsReq) ([]models.LeaveInfo, error) {
	userInfoRaw, err := database.DbConn.FindOne("employees", bson.D{{Key: "username", Value: employeeUsername}})
	if err != nil {
		return nil, err
	}

	var leavePipeline mongo.Pipeline
	if viewReq.Username == nil {
		var userInfo models.Employee
		if err := bson.Unmarshal(userInfoRaw, &userInfo); err != nil {
			log.SetPrefix("WARNING: ")
			log.Println("The decoding of employee from raw bson document failed!\nError:-\n\n", err)
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
		leavePipeline = AddUsernameFilterToDbPipeline(leavePipeline, peepsUsername)

	} else {
		leavePipeline = AddUsernameFilterToDbPipeline(leavePipeline, []string{*viewReq.Username})
	}
	leavePipeline = BuildLeaveFilterPipeline(leavePipeline, viewReq)
	leavePipeline = FlattenDbResult(leavePipeline, "$leaves")
	data, err := database.DbConn.Aggregate("leaves", leavePipeline)
	if err != nil {
		return nil, err
	}
	if data == nil {
		log.SetPrefix("WARNING: ")
		log.Println("No Leave entry found for given query")
		return nil, nil
	}

	leaves := database.ConvertRawBsonToLeaves(data)
	return leaves, nil
}
