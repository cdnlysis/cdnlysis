package pipeline

import (
	"io"
	"log"

	"cdnlysis/conf"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const LIMIT = 1000
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
	self.init(*marker)

	files := make(chan *LogFile)

	go func() {
		defer close(files)
		for _, file := range self.crawler.List.Contents {
			files <- &LogFile{file.Key, self.crawler.bucket}
			*marker = file.Key
		}
	}()

	return files
}

type LogFile struct {
	Path   string
	Bucket *s3.Bucket
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
