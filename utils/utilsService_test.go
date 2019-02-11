package utils

import (
	"testing"
	"time"
)

func TestFindMonthBirthdaySucess(t *testing.T) {
	repo := NewUtilsInMenRepository()
	service := NewUtilsService(repo)
	time, _ := time.Parse(time.RFC822, "01 May 15 10:00 UTC")
	births, err := service.FindMonthBirthday(time)
	if err != nil || births[0].Nome != "Bruno" {
		t.Error("Birts not listed")
	}
}

func TestFindMonthBirthdayWithFail(t *testing.T) {
	repo := NewUtilsInMenRepository()
	service := NewUtilsService(repo)
	time, _ := time.Parse(time.RFC822, "01 Jan 15 10:00 UTC")
	_, err := service.FindMonthBirthday(time)
	if err == nil {
		t.Error("Error no thrown")
	}
}
