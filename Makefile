export GO111MODULE=on

test:
	go test -v ./...

build:
	go build -o bin/oci-sd

clean:
	rm -rf ./bin