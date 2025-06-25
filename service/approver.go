package service

import (
	"LeafMS-BackEnd/database"
	"LeafMS-BackEnd/models"
	"LeafMS-BackEnd/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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

	data, err := database.DbConn.Aggregate("leaves", pipeline)
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
func ApproveLeave(leaveApplications models.ApproveLeaveReq) (*mongo.UpdateResult, error) {
	updatedResult, err := database.DbConn.UpdateOne("leaves", bson.D{
		{Key: "username", Value: leaveApplications.Username}, {
			Key: "leaves", Value: bson.D{{
				Key: "$elemMatch", Value: bson.D{{Key: "id", Value: leaveApplications.Leaves[0].Id}}}}}, //possible bug, why only matching for Leaves[0], why not for other IDs?
	}, bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "leaves.$.status", Value: leaveApplications.Leaves[0].Status},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return updatedResult, nil
}
