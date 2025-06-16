package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveType int
type LeaveStatus int

const (
	Sick LeaveType = iota
	Casual
	Paid
	Unpaid
	Maternity
	Paternity
)

const (
	Pending LeaveStatus = iota
	Approved
	Rejected
	Cancelled
)

type LeaveInfo struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EmployeeId primitive.ObjectID `bson:"employeeId" json:"employeeId"`
	Type       LeaveType          `bson:"type" json:"type"`
	StartDate  string             `bson:"startDate" json:"startDate"`
	EndDate    string             `bson:"endDate" json:"endDate"`
	Status     LeaveStatus        `bson:"status" json:"status"`
	Reason     string             `bson:"reason" json:"reason"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// type LeaveData struct {
// 	Id primitive.ObjectID `bson:"id" json:"id"`
// 	Type
// 	Start    string `bson:"startDate" json:"startDate"`
// 	End      string `bson:"endDate" json:"endDate"`
// 	Approved bool   `default:"false" bson:"approved" json:"approved"`
// }

type MetaLeaveInfo struct {
	Username string      `bson:"username" json:"username"`
	Approver string      `bson:"approver" json:"approver"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves"`
}
