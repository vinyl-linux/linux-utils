PREFIX ?= ""
BINDIR ?= "$(PREFIX)/bin"

BINARIES := $(BINDIR)/linux-utils
SCRIPTS := $(BINDIR)/useradd    \
	   $(BINDIR)/groupadd   \
	   $(BINDIR)/netctl


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

scripts/%:
	@echo "#!/bin/sh -e\n\n$(BINDIR)/linux-utils $* \"\$${@}\"" > $@

install: $(BINARIES) $(SCRIPTS)
