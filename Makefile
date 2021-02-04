GOCMD=go

all: build

build:
	$(GOCMD) build cli/groundhog.go

clean:
	rm -f groundhog
