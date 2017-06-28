package db

import (
	"fmt"
	"log"

	"github.com/adamluo159/gameAgent/utils"

	mgo "gopkg.in/mgo.v2"
)

var (
	// Session stores mongo session
	Session *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo

	MongoDBUrl string
)

// Connect connects to mongodb
func Connect() {
	err := utils.GetConfigValue("mongoIP", &MongoDBUrl)
	if err != nil {
		panic(fmt.Sprintf("read json file get mongodb ip fail %v\n", err))
	}
	log.Println(MongoDBUrl)
	s, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(fmt.Sprintf("Can't connect to mongoip:%s, go error %v\n", MongoDBUrl, err))
	}

	s.SetMode(mgo.Monotonic, true)
	fmt.Println("Connected to", MongoDBUrl)
	Session = s
}

func ReConnect() error {
	mongo, err := mgo.ParseURL(MongoDBUrl)
	s, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		fmt.Printf("Can't Reconnect to mongo, go error %v\n", err.Error())
		return err
	}
	s.SetMode(mgo.Monotonic, true)
	log.Println("ReConnected to", MongoDBUrl)

	Session.Close()
	Session = s
	Mongo = mongo
	return nil
}
