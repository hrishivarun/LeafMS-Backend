package models

type PostHoliday struct {
	Country string `bson:"country" json:"country"`
	Year    int    `bson:"year" json:"year"`
}
