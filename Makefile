build: deps
	go build cdnlysis

path:
	export GOPATH=`pwd`

deps: path
	go get launchpad.net/goamz/s3
	go get github.com/Simversity/gottp
	go get github.com/influxdb/influxdb/client
	go get github.com/HouzuoGuo/tiedot
