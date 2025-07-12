package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveType string
type LeaveStatus string

const (
	Sick      LeaveType = "Sick"
	Casual    LeaveType = "Casual"
	Paid      LeaveType = "Paid"
	Unpaid    LeaveType = "Unpaid"
	Maternity LeaveType = "Maternity"
	Paternity LeaveType = "Paternity"
)
const (
	Pending   LeaveStatus = "Pending"
	Approved  LeaveStatus = "Approved"
	Rejected  LeaveStatus = "Rejected"
	Cancelled LeaveStatus = "Cancelled"
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

//	type LeaveData struct {
//		Id primitive.ObjectID `bson:"id" json:"id"`
//		Type
//		Start    string `bson:"startDate" json:"startDate"`
//		End      string `bson:"endDate" json:"endDate"`
//		Approved bool   `default:"false" bson:"approved" json:"approved"`
//	}
type LeaveDoc struct {
	Username string      `bson:"username" json:"username" validate:"required"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves" validate:"required"`
}
