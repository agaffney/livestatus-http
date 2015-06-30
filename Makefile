GO ?= go
GOPATH := $(CURDIR)/_vendor:$(GOPATH)
PACKAGENAME=livestatus-http
PACKAGEVER=1.0
PACKAGEREL=1

all: build

build:
	$(GO) build

rpm: build
	fpm -s dir -t rpm --prefix /usr/bin --name $(PACKAGENAME) --version $(PACKAGEVER) --iteration $(PACKAGEREL) livestatus-http
