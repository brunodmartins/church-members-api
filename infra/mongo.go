package infra

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"gopkg.in/mgo.v2"
)

type mongoConnection struct{}

func NewMongoConnection() *mongoConnection {
	return &mongoConnection{}
}

func (mongoConn *mongoConnection) Connect() *mgo.Session {
	log.Println("Connecting to mongo")
	var mongoURI = os.Getenv("CHURCH_MEMBERS_DATABASE_URL")
	if mongoURI == "" {
		panic("CHURCH_MEMBERS_DATABASE_URL not defined")
	}
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
