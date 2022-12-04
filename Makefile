BINARY_NAME=orunmila
SOURCE=orunmila.go

.PHONY: all

all: dep test test_coverage vet lint build


build:
	GOARCH=amd64 GOOS=darwin  go build -o ${BINARY_NAME}-darwin      ${SOURCE}
	GOARCH=amd64 GOOS=linux   go build -o ${BINARY_NAME}-linux       ${SOURCE}
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows.exe ${SOURCE}
	GOARCH=amd64 GOOS=openbsd go build -o ${BINARY_NAME}-openbsd     ${SOURCE}

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	-@CGO_ENABLED=1 GO111MODULE=on go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all

run:
	./${BINARY_NAME}


cleandb:
	-rm ${BINARY_NAME}.db ${BINARY_NAME}.db-journal

clean:
	-go clean
	-rm ${BINARY_NAME}-darwin
	-rm ${BINARY_NAME}-linux
	-rm ${BINARY_NAME}-openbsd
	-rm ${BINARY_NAME}-windows.exe