package enum

type Role int
type Classification int

func (r Role) String() string {
	return []string{"ADMIN", "READ_ONLY"}[r]
}

func (c Classification) String() string {
	return []string{"Children", "Teen", "Young", "Adult"}[c]
}