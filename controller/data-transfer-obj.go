package controller

// ============================================================================
// ============================================================================
// DTO for viewing applications of people you're the leave approver of
// ============================================================================
// ============================================================================
type ViewApplications struct {
	ApproverName    string `bson:"approverName" json:"approverName"`
	IsLeaveAprroved *bool  `bson:"isLeaveAprroved" json:"isLeaveAprroved"`
}
