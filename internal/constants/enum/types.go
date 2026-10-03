package enum

type Classification int

func (c Classification) String() string {
	return []string{"Children", "Teen", "Young", "Adult"}[c]
}
