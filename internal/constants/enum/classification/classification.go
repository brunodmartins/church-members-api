package classification

import (
	"fmt"
	"github.com/brunodmartins/church-members-api/internal/constants/enum"
	"strings"
)

const (
	CHILDREN enum.Classification = iota
	TEEN
	YOUNG
	ADULT
)

var values = []string{"Children", "Teen", "Young", "Adult"}

func From(value string) (enum.Classification, error) {
	value = strings.ToUpper(value)
	for i := 0; i < len(values); i++ {
		if strings.ToUpper(values[i]) == value {
			return enum.Classification(i), nil
		}
	}
	return 0, fmt.Errorf("invalid classification: %s", value)
}
