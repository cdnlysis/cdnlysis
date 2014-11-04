package conf

import (
	utils "github.com/Simversity/gottp/utils"
	influxdb "github.com/influxdb/influxdb/client"
)

type Config struct {
	Influx influxdb.ClientConfig

	Engine struct {
		Verbose bool
		Threads int
	}

	Backends struct {
		Influx bool
		Mongo  bool
	}

	Mongo struct {
		Host       string
		Username   string
		Password   string
		Database   string
		Collection string
	}

	SyncProgress struct {
		Path string
	}

	S3 struct {
		Prefix    string
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
	}

	Logs struct {
		Prefix   string
		Location string
	}
}

const baseConfig = `;Sample Configuration File
[Backends]
Influx=true
Mongo=false

[Engine]
Verbose=true
Threads=10

[Influx]
Host="127.0.0.1:8086"
Username=root
Password=root
IsUDP=true
Database=cdn_logs

[Mongo]
Host="127.0.0.1:27017"
Username=""
Password=""
Database="cdn_logs"
Collection="cdn"

[SyncProgress]
Path="/tmp/cdn_sync_progress"

[S3]
Prefix = ""
AccessKey = ""
SecretKey = ""
Bucket = ""
Region = "us-east-1"

[Logs]
Prefix="cdn"
Location="/tmp/"`

func (self *Config) MakeConfig(configPath string) {
	utils.ReadConfig(baseConfig, self)
	if configPath != "" {
		utils.MakeConfig(configPath, self)
	}
}
