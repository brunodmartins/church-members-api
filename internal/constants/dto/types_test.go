package dto

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/brunodmartins/church-members-api/internal/constants/domain"
)

func TestUnmarshallDate(t *testing.T) {
	type TestUnmarshall struct {
		CustomDate Date `json:"customDate"`
	}
	result := &TestUnmarshall{}
	if err := json.Unmarshal([]byte(`{"customDate":"1995-05-10"}`), result); err != nil {
		assert.FailNow(t, err.Error())
	}
	assert.Equal(t, "1995-05-10", result.CustomDate.Format(time.DateOnly))
}

func TestMemberObservationRoundTrip(t *testing.T) {
	member := &domain.Member{
		Observation: "Needs follow-up from the pastoral team",
		Person: &domain.Person{
			Contact: &domain.Contact{},
			Address: &domain.Address{},
		},
		Religion: &domain.Religion{},
	}

	item := NewMemberItem(member)
	restored := item.ToMember()

	assert.Equal(t, member.Observation, restored.Observation)
}

func TestCreateMemberWithoutObservation(t *testing.T) {
	request := CreateMemberRequest{}
	err := json.Unmarshal([]byte(`{
		"attendsFridayWorship": true,
		"attendsSaturdayWorship": false,
		"attendsSundayWorship": true,
		"attendsSundaySchool": false
	}`), &request)
	assert.NoError(t, err)

	member := request.ToMember()
	assert.Equal(t, "", member.Observation)
	assert.NotNil(t, member.Person)
	assert.NotNil(t, member.Religion)
}
