package db

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

var (
	// Session stores mongo session
	Session *mgo.Session

	// Mongo stores the mongodb connection string information
	Mongo *mgo.DialInfo
)

const (
	// MongoDBUrl is the default mongodb url that will be used to connect to the
	// database.
	MongoDBUrl = "mongodb://192.168.1.252:27017/"
)

// Connect connects to mongodb
func Connect() {
	mongo, err := mgo.ParseURL(MongoDBUrl)
	s, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		panic(err.Error())
	}

	s.SetSafe(&mgo.Safe{})
	fmt.Println("Connected to", MongoDBUrl)
	fmt.Println("DialInfo", MongoDBUrl)
	Session = s
	Mongo = mongo
}
