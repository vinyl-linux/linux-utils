PREFIX ?= ""
BINDIR ?= "$(PREFIX)/bin"
CONFDIR ?= "$(PREFIX)/etc/vinyl/network.d"

BINARIES := $(BINDIR)/linux-utils
SCRIPTS := $(BINDIR)/useradd    \
	   $(BINDIR)/groupadd   \
	   $(BINDIR)/netctl

CONFIGS := $(CONFDIR)/eth0.toml.sample

.PHONY: default install

default: linux-utils

linux-utils:
	CGO_ENABLED=0 go build -o $@ ./bin/
	strip $@

$(BINDIR):
	mkdir -pv $@

$(CONFDIR):
	mkdir -pv $@

$(BINDIR)/linux-utils: linux-utils $(BINDIR)
	install -m 0750 -o root $< $@

$(BINDIR)/%: scripts/% $(BINDIR)
	install -m 0750 -o root $< $@

scripts/%: scripts
	@echo -e "#!/bin/sh -e\n\n$(BINDIR)/linux-utils $* \"\$${@}\"" > $@

scripts:
	mkdir -pv $@

$(CONFDIR)/eth0.toml.sample: $(CONFDIR)
	@echo -e 'interface = "eth0"\n\n[IPv4]\n dhcp = true\n enable = true' > $@

install: $(BINARIES) $(SCRIPTS) $(CONFIGS)
