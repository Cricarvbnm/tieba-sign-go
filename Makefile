SHELL = /bin/sh

app_name := tieba-sign
build_dir := build
app_build_path := $(build_dir)/$(app_name)

prefix := /usr/local
bindir := $(prefix)/bin
systemd_service_dir := $(prefix)/lib/systemd/user

INSTALL := install
INSTALL_DIR := $(INSTALL) -d
INSTALL_PROGRAM := $(INSTALL)
INSTALL_DATA := $(INSTALL) -m 644

src_make = $(MAKE) -C src app_name=$(app_name) build_dir=$(abspath $(build_dir))

SYSTEMD_UNITS := $(wildcard packaging/common/systemd/*)
systemd_targets := $(patsubst packaging/common/systemd/%, $(DESTDIR)$(systemd_service_dir)/%, $(SYSTEMD_UNITS))

.PHONY: all install clean FORCE

all: $(app_build_path)

FORCE:

$(app_build_path): FORCE
	@$(src_make)

install: all $(DESTDIR)$(bindir)/$(app_name) $(systemd_targets)

$(DESTDIR)$(bindir)/$(app_name): $(app_build_path)
	@$(INSTALL_DIR) "$(dir $@)"
	$(INSTALL_PROGRAM) "$<" "$@"

$(DESTDIR)$(systemd_service_dir)/%: packaging/common/systemd/%
	@$(INSTALL_DIR) "$(dir $@)"
	$(INSTALL_DATA) "$<" "$@"

clean:
	@$(src_make) clean
