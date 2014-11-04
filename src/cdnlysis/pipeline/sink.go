package pipeline

import (
	"cdnlysis/conf"
	"log"

	"cdnlysis/backends"

	influxdb "github.com/influxdb/influxdb/client"
	"labix.org/v2/mgo"
)

func AddToInflux(series *backends.InfluxRecord) {
	conn, err := influxdb.New(&conf.Settings.Influx)
	if err != nil {
		log.Println("Cannot connect to Influx", err)
		return
	}

	value := influxdb.Series(*series)

	if err := conn.WriteSeries([]*influxdb.Series{&value}); err != nil {
		log.Println("Cannot add to Influx", err)
		return
	}
}

func AddToMongo(records []interface{}) {

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
		return
	}

	conn := session.DB(conf.Settings.Mongo.Database)
	coll := conn.C(conf.Settings.Mongo.Collection)

	err = coll.Insert(records...)
	if err != nil {
		log.Println("Cannot Insert document in collection", err)
	}
}
