package models

import (
	"github.com/EDDYCJY/go-gin-example/pkg/setting"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/prometheus/common/log"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session
var chapterDB *mgo.Collection
var accountDB *mgo.Collection
var LogDB *mgo.Collection
var bookDB *mgo.Collection
var rankDB *mgo.Collection
var infoDB *mgo.Collection
var processDB *mgo.Collection
var EsDB *elasticsearch.Client

func SetEs() {
	cfg := elasticsearch.Config{
		Addresses: setting.EsSetting.Addrs,
		Username:  setting.EsSetting.Username,
		Password:  setting.EsSetting.Password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Create es connect failed:%n", err)
	}
	EsDB = es

}
func SetupMongo() {
	var err error
	mongo := &mgo.DialInfo{
		Addrs:    setting.MongoSetting.Addrs,
		Timeout:  setting.MongoSetting.Timeout,
		Database: setting.MongoSetting.Database,
		Username: setting.MongoSetting.Username,
		Password: setting.MongoSetting.Password,
	}
	session, err := mgo.DialWithInfo(mongo)
	if err != nil {
		log.Fatalf("CreateSession failed:%n", err)
	}

	//设置连接池的大小
	session.SetPoolLimit(setting.MongoSetting.PoolSize)

	if err != nil {
		log.Fatalf("set session pool size failed:%n", err)
	}

	chapterDB = session.DB("book").C(setting.MongoSetting.Chapter)
	//chapterDB := session.DB("book").C("chapter")
	accountDB = session.DB("book").C("account")
	bookDB = session.DB("book").C(setting.MongoSetting.Book)
	rankDB = session.DB("book").C("rank")
	LogDB = session.DB("book").C("log")
	infoDB = session.DB("book").C("info")
	processDB = session.DB("book").C("process")

}
func CloseMongo() {
	defer session.Close()
}
