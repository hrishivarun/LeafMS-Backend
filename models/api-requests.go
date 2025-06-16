package models

// ============================================================================
// ============================================================================
// DTO for viewing applications of people you're the leave approver of
// ============================================================================
// ============================================================================
type ViewApplications struct {
	ApproverName    string `bson:"approverName" json:"approverName"`
	IsLeaveAprroved *bool  `bson:"isLeaveAprroved" json:"isLeaveAprroved"`
}

type HolidaysFilter struct {
	Country string `bson:"country" json:"country"`
	Year    int    `bson:"year" json:"year"`
	Month   int    `bson:"month" json:"month"`
}
