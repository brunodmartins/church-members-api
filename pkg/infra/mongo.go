package infra

import (
	"crypto/tls"
	"log"
	"net"

	"gopkg.in/mgo.v2"
)

type mongoConnection struct{}

func NewMongoConnection() *mongoConnection {
	return &mongoConnection{}
}

func (mongoConn *mongoConnection) connect() *mgo.Session {
	log.Println("Connecting to mongo")
	var mongoURI = ""
	dialInfo, err := mgo.ParseURL(mongoURI)
	if err != nil {
		panic(err)
	}
	//Below part is similar to above.
	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, _ := mgo.DialWithInfo(dialInfo)

	session.SetMode(mgo.Monotonic, true)
	log.Println("Connected")
	return session
}
