package utils

import (
	"fmt"
	"time"
)

func ConvertDate(date time.Time) string {
	month := date.Month()
	day := date.Day()
	fmtMonth := fmt.Sprintf("%d", month)
	fmtDay := fmt.Sprintf("%d", day)
	if month < 10 {
		fmtMonth = fmt.Sprintf("0%d", month)
	}
	if day < 10 {
		fmtDay = fmt.Sprintf("0%d", day)
	}
	return fmt.Sprintf("%s-%s", fmtMonth, fmtDay)
}
