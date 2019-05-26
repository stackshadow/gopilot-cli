
deps:
	GOPATH=${HOME}/go:${PWD}/../gopilot-lib:${PWD} \
	go get -d -v ./src

build:
	GOPATH=${HOME}/go:${PWD}/../gopilot-lib:${PWD} \
	go build -o gpcli ./src

clean:
	rm gpcli