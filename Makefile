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
	go build -o $@ ./bin/

$(BINDIR):
	mkdir -pv $@

$(CONFDIR):
	mkdir -pv $@

$(BINDIR)/linux-utils: linux-utils $(BINDIR)
	install -m 0750 -o root $< $@

$(BINDIR)/%: scripts/% $(BINDIR)
	install -m 0750 -o root $< $@

scripts/%:
	@echo "#!/bin/sh -e\n\n$(BINDIR)/linux-utils $* \"\$${@}\"" > $@

$(CONFDIR)/eth0.toml.sample: $(CONFDIR)
	@echo -e 'interface = "eth0"\n\n[IPv4]\n dhcp = true\n enable = true' > $@

install: $(BINARIES) $(SCRIPTS) $(CONFIGS)
