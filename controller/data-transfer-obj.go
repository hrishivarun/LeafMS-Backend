package controller

// ============================================================================
// ============================================================================
// DTO for viewing applications of people you're the leave approver of
// ============================================================================
// ============================================================================
type ViewApplications struct {
	ApproverName    string `bson:"aprroverName" json:"aprroverName"`
	IsLeaveAprroved bool   `bson:"isLeaveAprroved" json:"isLeaveAprroved"`
}
