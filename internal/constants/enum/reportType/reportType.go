package reportType

const (
	MEMBER         = "members"
	LEGAL          = "legal"
	CLASSIFICATION = "classification"
	BIRTHDATE      = "birthdate"
	MARRIAGE       = "marriage"
)

var validReports = []string{MEMBER, LEGAL, CLASSIFICATION, BIRTHDATE, MARRIAGE}

func IsValidReport(name string) bool {
	for _, report := range validReports {
		if report == name {
			return true
		}
	}
	return false
}
