# Go parameters
GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOBENCH=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOCOVER=$(GOCMD) tool cover
GOGENERATE=$(GOCMD) generate
BINARY_NAME=go-test

all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

ci-test: test
	$(GOGET) -u github.com/jstemmer/go-junit-report
	cat test.out | go-junit-report > report.xml

test: build
	$(GOTEST) -v ./... --coverprofile=c.out | tee test.out
	$(GOCOVER) -html=c.out -o coverage.html
	$(GOCOVER) -func c.out | grep total

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	go mod download