package model

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
	"testing"
)


func TestID(t *testing.T) {
	t.Run("Test Get BSON", func(t *testing.T) {
		originalId := NewID()
		bsonValue, _ := originalId.GetBSON()
		objectId := bsonValue.(bson.ObjectId)
		assert.Equal(t, originalId.String(), objectId.Hex())
	})
}
