package main

import (
	utils "github.com/Simversity/gottp/utils"
	influxdb "github.com/influxdb/influxdb/client"
)

type config struct {
	Influx influxdb.ClientConfig

	S3 struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Region    string
	}

	Logs struct {
		Location string
	}
}

const baseConfig = `;Sample Configuration File
[Influx]
Host="127.0.0.1:8086"
Username=root
Password=root
IsUDP=true
Database=server_events

[S3]
AccessKey = ""
SecretKey = ""
Bucket = ""
Region = "us-east-1"

[Logs]
Location="/tmp/"`

func (self *config) MakeConfig(configPath string) {
	utils.ReadConfig(baseConfig, self)
	if configPath != "" {
		utils.MakeConfig(configPath, self)
	}
}
