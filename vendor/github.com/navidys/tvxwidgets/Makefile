TARGET := $(shell basename `pwd`)
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: check

check:
	@echo "running gofmt"
	@gofmt -l -w $(SRC)

	@echo "running golint"
	@for d in $$(go list ./... | grep -v /vendor/); do golint $${d}; done

	@echo "running go vet"
	@go vet ../$(TARGET)
	