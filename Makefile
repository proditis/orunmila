BINARY_NAME=orunmila
SOURCE=./...

.PHONY: all

all: dep test test_coverage vet build


build:
	GOARCH=amd64 GOOS=linux   go build -o ${BINARY_NAME}       ${SOURCE}

releases:
	GOARCH=amd64 GOOS=darwin  go build -o ${BINARY_NAME}-darwin      			 ${SOURCE}
	GOARCH=amd64 GOOS=linux   go build -o ${BINARY_NAME}-linux       			 ${SOURCE}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe 			 ${SOURCE}
	GOARCH=amd64 GOOS=openbsd go build -o ${BINARY_NAME}-openbsd     			 ${SOURCE}
	GOARCH=arm64 GOOS=linux   go build -o ${BINARY_NAME}-linux-arm64       ${SOURCE}
	GOARCH=arm64 GOOS=openbsd go build -o ${BINARY_NAME}-openbsd-arm64     ${SOURCE}

test:
	go test ${SOURCE}

test_coverage:
	go test ${SOURCE} -coverprofile=coverage.out

dep:
	-@CGO_ENABLED=1 GO111MODULE=on go mod download

vet:
	go vet

#lint:
#	golangci-lint run --enable-all

run:
	./${BINARY_NAME}


cleandb:
	-@rm ${BINARY_NAME}.db ${BINARY_NAME}.db-journal

clean:
	-go clean
	-@rm coverage.out
	-@rm ${BINARY_NAME}-darwin-amd64
	-@rm ${BINARY_NAME}-linux-amd64
	-@rm ${BINARY_NAME}-windows-amd64.exe
	-@rm ${BINARY_NAME}-openbsd-amd64
	-@rm ${BINARY_NAME}-linux-arm64
	-@rm ${BINARY_NAME}-openbsd-arm64
	-@rm ${BINARY_NAME}-darwin
	-@rm ${BINARY_NAME}-linux
	-@rm ${BINARY_NAME}-windows.exe
	-@rm ${BINARY_NAME}-openbsd
