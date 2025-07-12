package models

// type ViewApplicationsReq struct {
// 	EmployeeUsername *string `bson:"employeeUsername" json:"employeeUsername"`
// 	LeaveType        *string `bson:"leaveType" json:"leaveType"`
// 	Year             *int    `bson:"year" json:"year"`
// 	Month            *int    `bson:"month" json:"month"`
// 	Status           *string `bson:"status" json:"status"`
// }

type ResolveLeaveReq struct {
	Username *string     `bson:"username" json:"username"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves"`
}
