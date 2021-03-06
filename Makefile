
all: build

deps:
	go get github.com/go-sql-driver/mysql
	go get gopkg.in/yaml.v1

build:
	go build

build-linux:
	/usr/local/go/bin/linux_amd64/go build

ci: 
	cp default.yaml test.yaml

test:
	go test .

