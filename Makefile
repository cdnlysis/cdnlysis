path:
	export GOPATH=`pwd`

test: path
	echo ${GOPATH}

deps: path
	go get code.google.com/p/gcfg
	go get launchpad.net/goamz/s3
	go get github.com/Simversity/gottp
	go get github.com/influxdb/influxdb/client
	go get github.com/HouzuoGuo/tiedot

build: deps
	go build cdnlysis
