package model

import "gopkg.in/mgo.v2/bson"

//ID type
type ID bson.ObjectId

//ToString convert an ID in a string
func (i ID) String() string {
	return bson.ObjectId(i).Hex()
}

//StringToID convert a string to an ID
func StringToID(s string) ID {
	return ID(bson.ObjectIdHex(s))
}

//IsValidID check if is a valid ID
func IsValidID(s string) bool {
	return bson.IsObjectIdHex(s)
}

//NewID create a new id
func NewID() ID {
	return StringToID(bson.NewObjectId().Hex())
}
