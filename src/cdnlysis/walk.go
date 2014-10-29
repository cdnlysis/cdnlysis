package main

import (
	"io"
	"log"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const LIMIT = 1000
const STARTPOS = 0

func getRegion(cfg *config) aws.Region {
	return aws.Regions[cfg.S3.Region]
}

type s3iterator struct {
	bucket     *s3.Bucket
	prefix     string
	marker     string
	limit      int
	currentPos int
	crawler    *s3.ListResp
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

func (self *s3iterator) initCrawler() {
	self.currentPos = 0

	res, err := self.bucket.List(
		self.prefix, "", self.marker, self.limit,
	)

	if err != nil {
		log.Fatal(err)
	}

	if !res.IsTruncated {
		res.MaxKeys = len(res.Contents)
	}

	self.crawler = res
}

func (self *s3iterator) Next() *LogFile {
	if self.crawler == nil {
		self.initCrawler()
	}

	if self.currentPos >= self.crawler.MaxKeys &&
		self.crawler.MaxKeys == LIMIT {
		self.initCrawler()
	}

	key := self.crawler.Contents[self.currentPos]

	defer func(key s3.Key) {
		self.currentPos++
		self.marker = key.Key
	}(key)

	return &LogFile{key.Key, self.bucket}
}

func (self *s3iterator) End() bool {
	//return self.currentPos > 0
	if self.crawler != nil &&
		!self.crawler.IsTruncated &&
		self.currentPos >= self.crawler.MaxKeys {
		return true
	}

	return false
}

func getBucket(cfg *config) *s3.Bucket {
	auth := aws.Auth{
		AccessKey: cfg.S3.AccessKey,
		SecretKey: cfg.S3.SecretKey,
	}

	region := getRegion(cfg)

	connection := s3.New(auth, region)
	bucket := connection.Bucket(cfg.S3.Bucket)
	return bucket
}

func NewIterator(prefix string, cfg *config) *s3iterator {
	marker := ""
	bucket := getBucket(cfg)

	return &s3iterator{
		bucket, prefix, marker, LIMIT, STARTPOS, nil,
	}
}
