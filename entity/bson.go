package entity

import "gopkg.in/mgo.v2/bson"

// GetBSON implements bson.Getter.
func (i ID) GetBSON() (interface{}, error) {
	if i == "" {
		return "", nil
	}
	return bson.ObjectId(i), nil
}

// SetBSON implements bson.Setter.
func (i *ID) SetBSON(raw bson.Raw) error {
	decoded := new(string)
	bsonErr := raw.Unmarshal(decoded)
	if bsonErr == nil {
		*i = ID(bson.ObjectId(*decoded))
		return nil
	}
	return bsonErr
}