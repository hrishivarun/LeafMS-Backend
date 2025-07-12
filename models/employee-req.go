package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeaveApplication struct {
	Username string      `bson:"username" json:"username" validate:"required"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves" validate:"required,dive,required"`
}

type CancelLeavesReq struct {
	Username string               `bson:"user" json:"user" validate:"required"`
	LeaveIds []primitive.ObjectID `bson:"leaveIds" json:"leaveIds" validate:"required"`
}

type ViewApplicationsReq struct {
	Username  *string `bson:"employeeUsername" json:"employeeUsername"`
	LeaveType *string `bson:"leaveType" json:"leaveType"`
	Year      *int    `bson:"year" json:"year"`
	Month     *int    `bson:"month" json:"month"`
	Status    *string `bson:"status" json:"status"`
}

type HolidaysFilter struct {
	Country string `bson:"country" json:"country" validate:"required"`
	Year    int    `bson:"year" json:"year"`
	Month   int    `bson:"month" json:"month"`
}
