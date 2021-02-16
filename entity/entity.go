package entity

import "gopkg.in/mgo.v2/bson"

//ID type
type ID bson.ObjectId

//ToString convert an ID in a string
func (i ID) String() string {
	return bson.ObjectId(i).Hex()
}

// MarshalJSON will marshal ID to Json
func (i ID) MarshalJSON() ([]byte, error) {
	return bson.ObjectId(i).MarshalJSON()
}

// UnmarshalJSON will convert a string to an ID
func (i *ID) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]
	if bson.IsObjectIdHex(s) {
		*i = ID(bson.ObjectIdHex(s))
	}

	return nil
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
