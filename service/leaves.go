package service

import (
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

	result, err := dbConn.UpdateOne("leaves", bson.D{
		{Key: "username", Value: filteredLeaves.Username},
	}, bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "leaves", Value: bson.D{
				{Key: "$each", Value: filteredLeaves.Leaves},
			}},
		}},
	})

	if err != nil {
		log.Println("Encountered error while persisting applied leaves in Database. Err : ", err)
		return nil, err
	}
	return result, nil
}

// ============================================================================
// ============================================================================
// `view team's leaves`
// ============================================================================
// ============================================================================
func ViewTeamLeaveInfo(teamName string) ([]models.MetaLeaveInfo, error) {
	teamPeepsRaw, err := dbConn.Find("employees", bson.D{
		{Key: "team", Value: teamName}})

	if err != nil {
		return nil, err
	}
	teamPeeps := utils.ReturnEmployees(teamPeepsRaw)
	var peepsUsername []string
	for _, peep := range teamPeeps {
		peepsUsername = append(peepsUsername, peep.Username)
	}

	data, err := dbConn.Find("leaves", bson.D{
		{Key: "username", Value: bson.D{{Key: "$in", Value: peepsUsername}}}})

	if err != nil {
		return nil, err
	}

	leaves := utils.ReturnLeaves(data)
	return leaves, nil
}

// ============================================================================
// ============================================================================
// `view leave applications`
// ============================================================================
// ============================================================================
func ViewLeaveApplications(filter models.ViewApplications) ([]models.MetaLeaveInfo, error) {
	var pipeline mongo.Pipeline
	// Always filter by approver
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{
		{Key: "approver", Value: filter.ApproverName},
	}}})

	// If IsLeaveAprroved is provided, filter by it
	if filter.IsLeaveAprroved != nil {
		pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "leaves", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$leaves"},
					{Key: "as", Value: "leave"},
					{Key: "cond", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$$leave.approved", *filter.IsLeaveAprroved}},
					}},
				}},
			}},
		}}})
	}

	data, err := dbConn.Aggregate("leaves", pipeline)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, err
	}

	leaveApplications := utils.ReturnLeaves(data)
	return leaveApplications, nil
}

// ============================================================================
// ============================================================================
// `leaves approval`
// ============================================================================
// ============================================================================
func ApproveLeave(leaveApplications models.MetaLeaveInfo) (*mongo.UpdateResult, error) {
	updatedResult, err := dbConn.UpdateOne("leaves", bson.D{
		{Key: "username", Value: leaveApplications.Username}, {
			Key: "leaves", Value: bson.D{{
				Key: "$elemMatch", Value: bson.D{{Key: "id", Value: leaveApplications.Leaves[0].Id}}}}}, //possible bug, why only matching for Leaves[0], why not for other IDs?
	}, bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "leaves.$.approved", Value: leaveApplications.Leaves[0].Status},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return updatedResult, nil
}
