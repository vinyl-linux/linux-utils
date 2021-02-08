PREFIX ?= ""
BINDIR ?= "$(PREFIX)/bin"

BINARIES := $(BINDIR)/linux-utils
SCRIPTS := $(BINDIR)/useradd \
	   $(BINDIR)/groupadd


.PHONY: default install

default: linux-utils

linux-utils:
	go build -o $@ ./bin/

$(BINDIR):
	mkdir -pv $@

$(BINDIR)/linux-utils: linux-utils $(BINDIR)
	install -m 0750 -o root $< $@

$(BINDIR)/%: scripts/% $(BINDIR)
	install -m 0750 -o root $< $@

install: $(BINARIES) $(SCRIPTS)
