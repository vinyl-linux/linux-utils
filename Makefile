PREFIX ?= "linux-utils-x86_64"
BINDIR ?= "$(PREFIX)/bin"
CONFDIR ?= "$(PREFIX)/etc/vinyl/network.d"

BINARIES := $(BINDIR)/linux-utils
SCRIPTS := $(BINDIR)/useradd    \
	   $(BINDIR)/groupadd   \
	   $(BINDIR)/netctl     \
	   $(BINDIR)/wifi

CONFIGS := $(CONFDIR)/eth0.toml.sample

INSTALLERS := $(PREFIX)/Makefile

BUILT_ON := $(shell date --rfc-3339=seconds | sed 's/ /T/')
BUILT_BY := $(shell whoami)
BUILD_REF := $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

.PHONY: default package

default: linux-utils

linux-utils: pkg = "github.com/vinyl-linux/linux-utils/bin/cmd"
linux-utils: bin bin/*.go bin/**/*.go **/*.go
	CGO_ENABLED=0 go build -ldflags="-s -w -X $(pkg).Ref=$(BUILD_REF) -X $(pkg).BuildUser=$(BUILT_BY) -X $(pkg).BuiltOn=$(BUILT_ON)" -trimpath -o $@ ./$</

$(PREFIX):
	mkdir -pv $@

$(BINDIR):
	mkdir -pv $@

$(CONFDIR):
	mkdir -pv $@

$(BINDIR)/linux-utils: linux-utils | $(BINDIR)
	install -m 0755 $< $@

$(BINDIR)/%: scripts/% | $(BINDIR)
	install -m 0755 $< $@

scripts/%: | scripts
	@echo -e "#!/bin/sh -e\n\n$(BINDIR)/linux-utils $* \"\$${@}\"" > $@

scripts:
	mkdir -pv $@

$(CONFDIR)/eth0.toml.sample: $(CONFDIR)
	@echo -e 'interface = "eth0"\n\n[IPv4]\n dhcp = true\n enable = true' > $@

$(PREFIX)/%: | $(PREFIX)
	install -m 0640 $(subst $(PREFIX)/,,$@).package $@

package: $(BINARIES) $(SCRIPTS) $(CONFIGS) $(INSTALLERS)
