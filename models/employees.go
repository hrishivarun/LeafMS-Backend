package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id               primitive.ObjectID `bson:"id" json:"id"`
	Username         string             `bson:"username" json:"username" validate:"required"`
	Password         string             `bson:"password" json:"password" validate:"required"`
	FirstName        string             `bson:"firstName" json:"firstName" validate:"required"`
	MiddleName       string             `bson:"middleName,omitempty" json:"middleName,omitempty"`
	LastName         string             `bson:"lastName" json:"lastName"`
	Team             string             `bson:"team" json:"team" validate:"required"`
	Designation      string             `bson:"designation" json:"designation" validate:"required"`
	ApproverUserName string             `bson:"approver" json:"approver" validate:"required"`
	TotalLeaveCount  int                `bson:"leavesCount" json:"leavesCount" validate:"required"`
}

type LoginInfo struct {
	Username string `bson:"username" json:"username"`
	Status   int    `bson:"status" json:"status"`
	Token    string `bson:"token" json:"token"`
}

type LeaveCount struct {
	Sick        int `bson:"sick" json:"sick"`
	Casual      int `bson:"casual" json:"casual"`
	Paid        int `bson:"paid" json:"paid"`
	Unpaid      int `bson:"unpaid" json:"unpaid"`
	Maternity   int `bson:"maternity" json:"maternity"`
	Paternity   int `bson:"paternity" json:"paternity"`
	TotalLeaves int `bson:"totalCount" json:"totalCount"`
}

type LoginReq struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

// var defaultLeaveCount = LeaveCount{
// 	Sick:        10,
// 	Casual:      20,
// 	Paid:        0,
// 	Unpaid:      10,
// 	Maternity:   180,
// 	Paternity:   60,
// 	TotalLeaves: 365,
// }
