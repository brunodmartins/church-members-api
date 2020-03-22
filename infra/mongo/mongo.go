package mongo

import (
	"crypto/tls"
	"net"

	"github.com/BrunoDM2943/church-members-api/infra/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
)

type mongoConnection struct{}

func NewMongoConnection() *mongoConnection {
	return &mongoConnection{}
}

func (mongoConn *mongoConnection) Connect() *mgo.Session {
	log.Info("Connecting to mongo")
	var mongoURI = viper.GetString("mongo.url")

	if mongoURI == "" {
		log.Fatal("Mongo URL not defined")
	}

	var session *mgo.Session
	if config.IsProd() {
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
		session, _ = mgo.DialWithInfo(dialInfo)
	} else {
		session, _ = mgo.Dial(mongoURI)

	}

	session.SetMode(mgo.Monotonic, true)
	log.Info("Connected")
	return session
}
