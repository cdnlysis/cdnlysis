package main

import (
	utils "github.com/Simversity/gottp/utils"
	influxdb "github.com/influxdb/influxdb/client"
)

type config struct {
	Influx influxdb.ClientConfig

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
[Influx]
Host="127.0.0.1:8086"
Username=root
Password=root
IsUDP=true
Database=cdn_logs

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

func (self *config) MakeConfig(configPath string) {
	utils.ReadConfig(baseConfig, self)
	if configPath != "" {
		utils.MakeConfig(configPath, self)
	}
}
