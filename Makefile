.POSIX:
.SUFFIXES: .go

DESTDIR=

GO=go

OPS=

PREFIX=/usr

DEPS=proto/v1/v1.go proto/common.go

export GOPATH=$(shell pwd)/vendor/

all: options cinp-client cinp-server 

get:
	$(GO) get github.com/kesselborn/go-getopt
	mkdir -p $(GOPATH)/src/github.com/kakwa/
	[ -e $(GOPATH)/src/github.com/kakwa/cinp ] || ln -s `pwd` $(GOPATH)/src/github.com/kakwa/cinp

options:
	@echo "---"
	@echo "Compiler            = $(GO)"
	@echo "Compile options     = $(OPS)"
	@echo "GOPATH              = $(GOPATH)"
	@echo "Prefix              = $(PREFIX)"
	@echo "Install path server = $(DESTDIR)$(PREFIX)/bin/cinp-server"
	@echo "Install path client = $(DESTDIR)/sbin/cinp-client"
	@echo "---"

cinp-client: cinp-client.go $(DEPS)
	@echo "Building cinp-client"
	$(GO) $(OPS) build cinp-client.go

cinp-server: cinp-server.go $(DEPS)
	@echo "Building cinp-server"
	$(GO) $(OPS) build cinp-server.go

clean:
	@echo "Cleaning"
	@rm -f cinp-server
	@rm -f cinp-client

install: all
	@mkdir -p $(DESTDIR)$(PREFIX)/bin
	@mkdir -p $(DESTDIR)/sbin
	@cp cinp-client $(DESTDIR)/sbin/cinp-client
	@chmod 755 $(DESTDIR)/sbin/cinp-client
	@cp cinp-server $(DESTDIR)$(PREFIX)/bin/cinp-server
	@chmod 755 $(DESTDIR)$(PREFIX)/bin/cinp-server

uninstall:
	@echo removing executables from $(DESTDIR)$(PREFIX)/sbin
	@cd $(DESTDIR)$(PREFIX)/bin && rm -f cinp-server
	@cd $(DESTDIR)/sbin && rm -f cinp-client

.PHONY: all options clean install uninstall get
