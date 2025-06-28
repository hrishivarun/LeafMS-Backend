package models

type ViewApplications struct {
	ApproverName    string `bson:"approverName" json:"approverName"`
	IsLeaveAprroved *bool  `bson:"isLeaveAprroved" json:"isLeaveAprroved"`
}

type ResolveLeaveReq struct {
	Username string      `bson:"username" json:"username"`
	Leaves   []LeaveInfo `bson:"leaves" json:"leaves"`
}
