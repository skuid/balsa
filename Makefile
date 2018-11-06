REPO=balsa

.PHONY: clean

all: build

test:
	@go test -cover ./...

build:
	@go build

clean:
	rm ./$(REPO)
