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
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type      LeaveType          `bson:"type" json:"type" validate:"required"`
	StartDate time.Time          `bson:"startDate" json:"startDate" validate:"required"`
	EndDate   time.Time          `bson:"endDate" json:"endDate" validate:"required"`
	Status    LeaveStatus        `bson:"status" json:"status" validate:"required"`
	Reason    string             `bson:"reason" json:"reason"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

// type LeaveData struct {
// 	Id primitive.ObjectID `bson:"id" json:"id"`
// 	Type
// 	Start    string `bson:"startDate" json:"startDate"`
// 	End      string `bson:"endDate" json:"endDate"`
// 	Approved bool   `default:"false" bson:"approved" json:"approved"`
// }
