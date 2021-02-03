PREFIX ?= ""
BINDIR ?= "$(PREFIX)/bin"

.PHONY: default install

default: linux-utils

linux-utils:
	go build -o $@ ./bin/

$(BINDIR):
	mkdir -pv $@

$(BINDIR)/linux-utils: linux-utils $(BINDIR)
	install -m 0750 -o root $< $@

$(BINDIR)/useradd: scripts/useradd $(BINDIR)
	install -m 0750 -o root $< $@

$(BINDIR)/groupadd: scripts/groupadd $(BINDIR)
	install -m 0750 -o root $< $@

install: $(BINDIR)/linux-utils $(BINDIR)/useradd $(BINDIR)/groupadd
