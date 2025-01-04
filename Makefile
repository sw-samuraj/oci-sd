export GO111MODULE=on

test:
	go test -v -cover ./...

build:
	go build -mod vendor -o /bin/oci-sd

clean:
	rm -rf ./bin