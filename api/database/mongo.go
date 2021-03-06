package database

import (
	"gAPIManagement/api/utils"
	"gopkg.in/mgo.v2"
)

type MongoPool struct {
	BaseSession *mgo.Session
	Queue chan int
	URL string
	Open int
}
var MongoDBPool MongoPool


func GetSessionAndDB(db string) (*mgo.Session, *mgo.Database) {
	session := GetSession()
	dbConn := GetDB(session, db)

	return session, dbConn
}

func (mp *MongoPool) New() error {
	var err error
	maxPool := 50
	mp.Queue = make(chan int, maxPool)
	for i := 0; i < maxPool; i = i + 1 {
		mp.Queue <- 1
	}
	mp.Open = 0
	mp.BaseSession, err = mgo.Dial(mp.URL)
	
	return err
}

func (mp *MongoPool) Session() *mgo.Session {
	defer utils.PreventCrash()

	<- mp.Queue
	if mp.BaseSession == nil {
		mp.BaseSession, _ = mgo.Dial(mp.URL)
	}
	mp.Open++
	return mp.BaseSession.Clone()
}
func (mp *MongoPool) Close(session *mgo.Session) {
	session.Close()

	mp.Queue <- 1
	mp.Open--
}


func ConnectToMongo(host string) error{
	var err error
	MongoDBPool.URL = host
	err = MongoDBPool.New()

	if err != nil {
		utils.LogMessage("error connecting to mongo on " + host + ". Err: " + err.Error(), utils.ErrorLogType)
	} else {
		utils.LogMessage("connected to mongo on " + host, utils.InfoLogType)
	}

	return err
}

func GetSession() *mgo.Session {	
	return MongoDBPool.Session()
}

func GetDB(session *mgo.Session, db string) *mgo.Database {
	return session.DB(db)
}

func Query(q interface{}) error {

	return nil
}
