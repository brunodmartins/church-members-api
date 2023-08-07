package dto

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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
