package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthorizationStatus int

const (
	StatusRequested AuthorizationStatus = iota
	StatusUnauthorized
	StatusForbidden
	StatusNotFound
	StatusInternalServerError
)

// ============================================================================
// ============================================================================
// DTO for viewing applications of people you're the leave approver of
// ============================================================================
// ============================================================================
type ViewApplications struct {
	ApproverName    string `bson:"approverName" json:"approverName"`
	IsLeaveAprroved *bool  `bson:"isLeaveAprroved" json:"isLeaveAprroved"`
}

type CancelLeavesReq struct {
	Username string               `bson:"user" json:"user" validate:"required"`
	LeaveIds []primitive.ObjectID `bson:"leaveIds" json:"leaveIds"`
}

type ApproveLeaveReq struct {
	Username string      `bson:"username" json:"username"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves"`
}

type HolidaysFilter struct {
	Country string `bson:"country" json:"country"`
	Year    int    `bson:"year" json:"year"`
	Month   int    `bson:"month" json:"month"`
}

type PostHoliday struct {
	Country string `bson:"country" json:"country"`
	Year    int    `bson:"year" json:"year"`
}
