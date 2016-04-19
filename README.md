![Build Status](https://travis-ci.org/cdnlysis/cdnlysis.svg?branch=dev)

cdnlysis
========

CDNlysis syncs Amazon Cloudfront log entries from S3 bucket and streams them to multiple database backends. As of now it writes to InfluxDB & MongoDB (turned off by default).
You can later use Influx Querying API to do anaylsis or use Grafana to generate awesome graphs.

You can use this for:
* Understanding how the Bandwidth is being used.
* Finding out the most popular and most downloadable content.
* Generate trends for your most popular Videos, Audios, Slides etc.
* Understand geographical behaviour of the Requests.
* Amount of Bytes transferred to & fro the Cloudfront distributions.
* Find out the most profitable referrer from where your content is being accessed.
etc.

# Features
* CDNlysis resumes a sync operation from the last processed file.
* CDN logs are saved as gzip files in S3, CDNlysis downloads, unzips, parses and feeds this to InfluxDB.
* It iteratively walks over your historical logFiles too.
* CDNlysis converts the time provided in the log entries as the actual time of saved records. This ensures consistency in expctations and data delivered.
* CDNlysis works not only over Amazon Cloudfront but anything that complies with W3C Extended Log Format (http://www.w3.org/TR/WD-logfile.html).
* CDNlysis can write to either InfluxDB, MongoDB or both.

# Usage
Using CDNlysis is dirt easy, just make sure you have latest Go installed and follow these steps:
Create a simple file called run.go
```
// +build !appengine

package main

import (
	cdnlysis "gopkg.in/cdnlysis/cdnlysis.v1"
	conf "gopkg.in/cdnlysis/cdnlysis.v1/conf"
	db "gopkg.in/cdnlysis/cdnlysis.v1/db"
)

func main() {
	cdnlysis.Setup()
	
	marker := db.LastMarker(conf.Settings.S3.Prefix)
	//Only use this to resume from the last saved anchor.
	
	record_chan := make(chan *cdnlysis.LogRecord)
	cdnlysis.Start(&marker, record_chan)
}
```

```
go get github.com/cdnlysis/cdnlysis
go run run.go
```

# Configuration

To generate a Binary, just run go build in the directoty where your run.go is

```
piyush:cdnlysis  λ go build .

piyush:cdnlysis  λ ls
cdnlysis  sample.go
```

Use the newly created binary to generate help

```
piyush:cdnlysis [dev] $ ./cdnlysis
  -config="": Config [.ini format] file to Load the configurations from
  -prefix="": Directory prefix to process the logs for
```

-config should be a path to a valid configuration file which can have 4 sections.
 * Backends: Quick access switch for the enabled Backends.
 * Influx: Influx configuration.
 * Mongo: Mongo configuration.
 * SyncProgress: CDNlysis' internal database. The path must be changed to something more permanent. This will help you prevent from feeding redundant entries to your Influx database.
 * S3: AWS configurations. Prefix is the prefix/directory in which logs should be searched.

A sample configuration file looks like this:

```
[Backends]
Influx=true
Mongo=false

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
```

Optionally, you can also provide -prefix as a Command Line Argument, by which you can run multiple sync jobs corresponding to each directory.

# Execute
Just provide the config path and an optional prefix. Sit back and let the tool do the rest.

```
./cdnlysis -config=/Users/piyush/dev/conf/cdn-analysis.cfg -prefix=plain
```

# Screenshots
![Screenshot1](https://cloud.githubusercontent.com/assets/580782/4833122/f1baa26a-5fa1-11e4-919e-261f46cec2b0.png)
![Screenshot2](https://cloud.githubusercontent.com/assets/580782/4833123/f1bb5002-5fa1-11e4-910c-35a4845843e0.png)
![Screenshot3](https://cloud.githubusercontent.com/assets/580782/4833124/f1eff384-5fa1-11e4-99a3-b35876566ccd.png)
