package member

type EventType string

const (
	EventMemberCreated EventType = "MEMBER_CREATED"
	EventMemberUpdated EventType = "MEMBER_UPDATED"
	EventMemberRetired EventType = "MEMBER_RETIRED"
)

type Event struct {
	MemberID  string    `json:"memberId"`
	EventType EventType `json:"eventType"`
}
