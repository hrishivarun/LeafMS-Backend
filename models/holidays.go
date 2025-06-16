package models

//structs for fetching public holidays
type Meta struct {
	Code int `bson:"code" json:"code"`
}
type Country struct {
	Id   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}
type Datetime struct {
	Day   int `bson:"day" json:"day"`
	Month int `bson:"month" json:"month"`
	Year  int `bson:"year" json:"year"`
}
type Date struct {
	Iso      string   `bson:"iso" json:"iso"`
	Datetime Datetime `bson:"datetime" json:"datetime"`
}
type Holiday struct {
	Name         string   `bson:"name" json:"name"`
	Description  string   `bson:"description" json:"description"`
	Country      Country  `bson:"country" json:"country"`
	Date         Date     `bson:"date" json:"date"`
	Type         []string `bson:"type" json:"type"`
	PrimaryType  string   `bson:"primary_type" json:"primary_type"`
	CanonicalUrl string   `bson:"canonical_url" json:"canonical_url"`
	UrlId        string   `bson:"urlid" json:"urlid"`
	Locations    string   `bson:"locations" json:"locations"`
	States       string   `bson:"states" json:"states"`
}
type HolidayResponse struct {
	Holidays []Holiday `bson:"holidays" json:"holidays"`
}
type HolidayApiResponse struct {
	Meta     Meta            `bson:"meta" json:"meta"`
	Response HolidayResponse `bson:"response" json:"response"`
}

func (date1 Datetime) IsGreaterThanOrEquals(date2 Datetime) bool {
	if date1.Year > date2.Year {
		return true
	} else if date1.Year < date2.Year {
		return false
	}

	if date1.Month > date2.Month {
		return true
	} else if date1.Month < date2.Month {
		return false
	}
	if date1.Day > date2.Day {
		return true
	} else if date1.Day < date2.Day {
		return false
	}
	return true
}
