package pipeline

import (
	"github.com/meson10/cdnlysis/conf"
	"log"

	"github.com/meson10/cdnlysis/backends"

	influxdb "github.com/influxdb/influxdb/client"
	"labix.org/v2/mgo"
)

var cachedInfluxConn *influxdb.Client

func makeInfluxSession() error {
	conn, err := influxdb.New(&conf.Settings.Influx)
	if err != nil {
		log.Println("Cannot connect to Influx", err)
		return err
	}

	cachedInfluxConn = conn
	return nil
}

func RefreshInfluxSession() error {
	log.Println("Refreshing Influx Connection")
	err := makeInfluxSession()
	return err
}

func AddToInflux(series *backends.InfluxRecord) {
	if cachedInfluxConn == nil {
		makeInfluxSession()
	}

	conn := cachedInfluxConn
	value := influxdb.Series(*series)

	if err := conn.WriteSeries([]*influxdb.Series{&value}); err != nil {
		log.Println("Cannot add to Influx", err)
		return
	}
}

var cachedSession *mgo.Session

func makeSession() error {
	info := mgo.DialInfo{
		Addrs:    []string{conf.Settings.Mongo.Host},
		Database: conf.Settings.Mongo.Database,
		Direct:   true,
		Username: conf.Settings.Mongo.Username,
		Password: conf.Settings.Mongo.Password,
	}

	var err error
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		log.Println("Cannot connect to Mongo:", err)
		return err
	}

	cachedSession = session
	return nil
}

func RefreshMongoSession() error {
	log.Println("Refreshing Mongo Connection")
	err := makeSession()
	if err == nil {
		cachedSession.Refresh()
	}
	return err
}

func AddToMongo(records []interface{}) {
	if cachedSession == nil {
		makeSession()
	}

	session := cachedSession.Copy()
	defer session.Close()

	conn := session.DB(conf.Settings.Mongo.Database)
	coll := conn.C(conf.Settings.Mongo.Collection)

	err := coll.Insert(records...)
	if err != nil {
		log.Println("Cannot Insert document in collection", err)
	}
}
