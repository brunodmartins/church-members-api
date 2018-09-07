package utils

import (
	"testing"
	"time"
)

func TestRepoFindMonthBirthdaySucess(t *testing.T) {
	repo := NewUtilsInMenRepository()
	time, _ := time.Parse(time.RFC822, "01 May 15 10:00 UTC")
	births, err := repo.FindMonthBirthday(time)
	if err != nil || births[0].Nome != "Bruno" {
		t.Error("Birts not listed")
	}
}

func TestRepoFindMonthBirthdayWithFail(t *testing.T) {
	repo := NewUtilsInMenRepository()
	time, _ := time.Parse(time.RFC822, "01 Jan 15 10:00 UTC")
	_, err := repo.FindMonthBirthday(time)
	if err == nil {
		t.Error("Error no thrown")
	}
}
