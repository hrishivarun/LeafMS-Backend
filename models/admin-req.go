package models

type EmployeesCrudReq struct {
	Usernames []string `bson:"usernames" json:"usernames" validate:"required,dive,required"`
}

type PostHoliday struct {
	Country string `bson:"country" json:"country" validate:"required"`
	Year    int    `bson:"year" json:"year" validate:"required"`
}
