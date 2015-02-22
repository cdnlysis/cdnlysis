package cdnlysis

import (
	"io"
	"log"
	"strconv"

	"github.com/meson10/cdnlysis/conf"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const LIMIT = 500
const STARTPOS = 0

func getRegion() aws.Region {
	return aws.Regions[conf.Settings.S3.Region]
}

func getBucket() *s3.Bucket {
	auth := aws.Auth{
		AccessKey: conf.Settings.S3.AccessKey,
		SecretKey: conf.Settings.S3.SecretKey,
	}

	region := getRegion()

	connection := s3.New(auth, region)
	bucket := connection.Bucket(conf.Settings.S3.Bucket)
	return bucket
}

type s3crawler struct {
	List       *s3.ListResp
	bucket     *s3.Bucket
	currentPos int
}

type Iterator struct {
	crawler     *s3crawler
	IsTruncated bool
}

func (self *Iterator) init(marker string) {
	bucket := getBucket()
	res, err := bucket.List(conf.Settings.S3.Prefix, "", marker, LIMIT)
	if err != nil {
		log.Fatal(err)
	}

	self.IsTruncated = res.IsTruncated
	self.crawler = &s3crawler{res, bucket, 0}
}

func (self *Iterator) Produce(marker *string) <-chan *LogFile {
	log.Println("[Producer] Fetching", LIMIT, "objects")
	self.init(*marker)

	files := make(chan *LogFile)
	count := len(self.crawler.List.Contents)

	log.Println("[Producer] Dispatching", count, "objects")

	go func(count int) {
		defer close(files)
		for ix, file := range self.crawler.List.Contents {
			logfile := LogFile{file.Key, self.crawler.bucket, ix, count}
			log.Println(logfile.LogIdent(), "[Producer] Dispatching")
			files <- &logfile
			*marker = file.Key
		}
	}(count)

	return files
}

type LogFile struct {
	Path      string
	Bucket    *s3.Bucket
	Sequence  int
	BatchSize int
}

func (self *LogFile) LogIdent() string {
	return strconv.Itoa(self.Sequence) + "/" + strconv.Itoa(self.BatchSize) + " " + self.Path
}

func (self *LogFile) Get() (data []byte, err error) {
	return self.Bucket.Get(self.Path)
}

func (self *LogFile) GetReader() (io.ReadCloser, error) {
	return self.Bucket.GetReader(self.Path)
}

func NewIter(marker string) *Iterator {
	return &Iterator{nil, false}
}
