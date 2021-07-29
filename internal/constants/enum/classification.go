package enum

import (
	"fmt"
	"strings"
)

type Classification int

const (
	CHILDREN Classification = iota + 1
	TEEN
	YOUNG
	ADULT
)


var classifications = []string{"Children", "Teen", "Young", "Adult"}

func (c Classification) String() string {
	return classifications[c - 1]
}

func (Classification) From(value string) (Classification, error) {
	value = strings.ToUpper(value)
	for i := 0; i < len(classifications); i++ {
		if strings.ToUpper(classifications[i]) == value {
			return Classification(i + 1), nil
		}
	}
	return 0, fmt.Errorf("invalid classification: %s", value)
}