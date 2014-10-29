cdnlysis
========

CDNlysis reads your log entries from your S3 bucket and sends them your Influx Server. You can later use Influx Querying API to do anaylsis or use Grafana to generate awesome graphs.

You can use this for:
* Understanding how the Bandwidth is being used.
* Finding out the most downloadable content
* Generate trends for your most popular S3 content.
* Understand geographical behaviour of your content access.
* Amount of Bytes transferred to & fro the Cloudfront distributions.
* Find out the most profitable referrer from where your content is being accessed.
etc.

# Features
* CDNlysis resumes a sync operation from the last processed file.
* CDN logs are saved as gzip files in S3, CDNlysis downloads, unzips, parses and feeds this to InfluxDB.
* It iteratively walks over your historical logFiles too.
* CDNlysis converts the time provided in the log entries as the actual time of saved records. This ensures consistency in expctations and data delivered.
* CDNlysis works not only over Amazon Cloudfront but anything that complies with W3C Extended Log Format (http://www.w3.org/TR/WD-logfile.html).

# Usage
Installing CDNlysis is dirty easy, just make sure you have latest Go installed and follow these steps:
In the checkout directory

```
make deps
```
Should install all the dependencies for you.

```
make build
```
should generate a binary for you to execute.

# Configuration

Use the newly created binary to generate help

```
piyush:cdnlysis [dev] $ ./cdnlysis
  -config="": Config [.ini format] file to Load the configurations from
  -prefix="": Directory prefix to process the logs for
```

-config should be a path to a valid configuration file which can have 4 sections.
 * Influx has the all the information for your Influx deployment
 * SyncProgress is CDNlysis' internal database. The path must be changed to something more permanent. This will help you prevent from feeding redundant entries to your Influx database.
 * S3 keeps all your AWS configurations. Prefix is the prefix/directory in which logs should be searched

A sample configuration file looks like this:

```
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
```

Optionally, you can also provide -prefix as a Command Line Argument, by which you can run multiple sync jobs corresponding to each directory.

# Execute
Just provide the config path and an optional prefix. Sit back and let the tool do the rest.

```
./cdnlysis -config=/Users/piyush/dev/conf/cdn-analysis.cfg -prefix=plain
```

# Screenshots
![Screenshot 1](https://www.dropbox.com/s/ncf8e25noenfy2f/Screen%20Shot%202014-10-30%20at%2000.49.08.png?dl=0)
