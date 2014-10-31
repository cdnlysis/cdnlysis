package main

import (
	"io"
	"log"

	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

const LIMIT = 50
const STARTPOS = 0

func getRegion(cfg *config) aws.Region {
	return aws.Regions[cfg.S3.Region]
}

type s3crawler struct {
	List       *s3.ListResp
	bucket     *s3.Bucket
	currentPos int
}

type s3iterator struct {
	prefix      string
	cfg         *config
	crawler     *s3crawler
	IsTruncated bool
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

func (self *s3iterator) initCrawler(marker string) {
	bucket := getBucket(self.cfg)
	res, err := bucket.List(self.prefix, "", marker, LIMIT)
	if err != nil {
		log.Fatal(err)
	}

	self.IsTruncated = res.IsTruncated

	if !res.IsTruncated {
		res.MaxKeys = len(res.Contents)
	}

	self.crawler = &s3crawler{res, bucket, 0}
}

func (self *s3iterator) Next() *LogFile {
	key := self.crawler.List.Contents[self.crawler.currentPos]

	defer func(key s3.Key) {
		self.crawler.currentPos++
	}(key)

	return &LogFile{key.Key, self.crawler.bucket}
}

func (self *s3iterator) End() bool {
	if self.crawler.currentPos >= self.crawler.List.MaxKeys {
		return true
	} else if self.crawler.List.MaxKeys == 0 {
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

func NewIterator(prefix string, marker string, cfg *config) *s3iterator {
	iter := &s3iterator{prefix, cfg, nil, false}
	iter.initCrawler(marker)
	return iter
}
