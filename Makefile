
all: build

deps:
	go get github.com/go-sql-driver/mysql
	go get gopkg.in/yaml.v1
	go get github.com/stretchr/testify/mock

build:
	go build

build-linux:
	/usr/local/go/bin/linux_amd64/go build

test:
	go test .

