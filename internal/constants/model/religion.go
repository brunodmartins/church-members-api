package model

import (
	"time"
)

//Religion struct
type Religion struct {
	FathersReligion   string     `json:"fathersReligion,omitempty"`
	BaptismPlace      string     `json:"baptismPlace"`
	LearnedGospelAge  int        `json:"learnedGospelAge"`
	AcceptedJesus     bool       `json:"acceptedJesus"`
	Baptized          bool       `json:"baptized"`
	CatholicBaptized  bool       `json:"catholicBaptized"`
	KnowsTithe        bool       `json:"knowsTithe"`
	AgreesTithe       bool       `json:"agreesTithe"`
	Tithe             bool       `json:"tithe"`
	AcceptedJesusDate *time.Time `json:"acceptedJesusDate"`
	BaptismDate       *time.Time `json:"baptismDate"`
}
