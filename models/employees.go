package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	Id                 primitive.ObjectID `bson:"id" json:"id"`
	Username           string             `bson:"username" json:"username" validate:"required"`
	Password           string             `bson:"password" json:"password" validate:"required"`
	FirstName          string             `bson:"firstName" json:"firstName" validate:"required"`
	MiddleName         string             `bson:"middleName,omitempty" json:"middleName,omitempty"`
	LastName           string             `bson:"lastName" json:"lastName"`
	Team               string             `bson:"team" json:"team" validate:"required"`
	Designation        string             `bson:"designation" json:"designation" validate:"required"`
	ApproverUserName   string             `bson:"approver" json:"approver" validate:"required"`
	LeavesCapacityLeft map[string]int     `bson:"leavesCapacityLeft" json:"leavesCapacityLeft" validate:"required"`
}

type LoginInfo struct {
	Username string `bson:"username" json:"username"`
	Status   int    `bson:"status" json:"status"`
	Token    string `bson:"token" json:"token"`
}

type LoginReq struct {
	Username string `bson:"username" json:"username" validate:"required"`
	Password string `bson:"password" json:"password" validate:"required"`
}

var DefaultLeavesCapacity = map[string]int{
	"Sick":      10,
	"Casual":    24,
	"Paid":      30,
	"Unpaid":    10,
	"Maternity": 180,
	"Paternity": 60,
}

func CopyDefaultLeavesCapacity() map[string]int {
	copy := make(map[string]int, len(DefaultLeavesCapacity))
	for k, v := range DefaultLeavesCapacity {
		copy[k] = v
	}
	return copy
}
