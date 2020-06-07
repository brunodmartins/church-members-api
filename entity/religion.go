package entity

import (
	"time"
)

//Religion struct
type Religion struct {
	FathersReligion   string    `json:"fathersReligion,omitempty"`
	BaptismPlace      string    `json:"baptismPlace" bson:"baptismPlace"`
	LearnedGospelAge  int       `json:"learnedGospelAge" bson:"learnedGospelAge"`
	AcceptedJesus     bool      `json:"acceptedJesus" bson:"acceptedJesus"`
	Baptized          bool      `json:"baptized"`
	CatholicBaptized  bool      `json:"catholicBaptized" bson:"catholicBaptized"`
	KnowsTithe        bool      `json:"knowsTithe" bson:"knowsTithe"`
	AgreesTithe       bool      `json:"agreesTithe" bson:"agreesTithe"`
	Tithe             bool      `json:"tithe"`
	AcceptedJesusDate time.Time `json:"acceptedJesusDate" bson:"acceptedJesusDate"`
	BaptismDate       time.Time `json:"baptismDate" bson:"baptismDate"`
}
