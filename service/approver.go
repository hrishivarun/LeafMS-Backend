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
// `view leave applications`
// ============================================================================
// ============================================================================
func ViewLeaveApplications(req models.ViewApplicationsReq, approverUsername string) ([]models.LeaveInfo, error) {
	var leaveApplicationPipeline mongo.Pipeline
	if req.Username != nil {
		leaveApplicationPipeline = AddUsernameFilterToDbPipeline(leaveApplicationPipeline, []string{*req.Username})
	} else {
		leaveApplicationPipeline = AddApproverFilterToDbPipeline(leaveApplicationPipeline, approverUsername)
	}
	leaveApplicationPipeline = BuildLeaveFilterPipeline(leaveApplicationPipeline, req)
	leaveApplicationPipeline = FlattenDbResult(leaveApplicationPipeline, "$leaves")

	data, err := database.DbConn.Aggregate("leaves", leaveApplicationPipeline)
	if err != nil {
		return nil, err
	}
	if data == nil {
		log.SetPrefix("WARNING: ")
		log.Println("No Leave entry found for given query")
		return nil, nil
	}

	leaveApplications := database.ConvertRawBsonToLeaves(data)
	return leaveApplications, nil
}

// ============================================================================
// ============================================================================
// `resolve leaves`
// ============================================================================
// ============================================================================
func ResolveLeave(req models.ResolveLeaveReq) ([]*mongo.UpdateResult, error) {
	var results []*mongo.UpdateResult

	for _, leave := range req.Leaves {
		filter := bson.D{
			{Key: "username", Value: req.Username},
			{Key: "leaves", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "id", Value: leave.Id},
				}},
			}},
		}
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "leaves.$.status", Value: leave.Status},
			}},
		}

		result, err := database.DbConn.UpdateOne("leaves", filter, update, nil)
		if err != nil {
			return results, err // return what was done so far and the error
		}
		results = append(results, result)
	}

	return results, nil
}
