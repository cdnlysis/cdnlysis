package main

import (
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
	prefix     string
	marker     string
	limit      int
	currentPos int
	crawler    *s3.ListResp
	cfg        *config
}

func (self *s3iterator) initCrawler() {
	self.currentPos = 0

	auth := aws.Auth{
		AccessKey: self.cfg.S3.AccessKey,
		SecretKey: self.cfg.S3.SecretKey,
	}

	region := getRegion(self.cfg)

	connection := s3.New(auth, region)
	bucket := connection.Bucket(self.cfg.S3.Bucket)
	res, err := bucket.List(self.prefix, "", self.marker, self.limit)
	if err != nil {
		log.Fatal(err)
	}

	if !res.IsTruncated {
		res.MaxKeys = len(res.Contents)
	}

	self.crawler = res
}

func (self *s3iterator) Next() *s3.Key {
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

	return &key
}

func (self *s3iterator) End() bool {
	if self.crawler != nil &&
		!self.crawler.IsTruncated &&
		self.currentPos >= self.crawler.MaxKeys {
		return true
	}

	return false
}

func NewIterator(prefix string, cfg *config) *s3iterator {
	marker := ""

	return &s3iterator{
		prefix, marker, LIMIT, STARTPOS, nil, cfg,
	}
}
