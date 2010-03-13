# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

ifneq ($(strip $(shell which gotmake)),)
all: packages Makefile

Makefile: scripts/make.header scripts/mkmake $(wildcard */*/.go) $(wildcard */*.go)
	./scripts/mkmake
else
all: packages
endif

test: install testit

install: installpkgs


include $(GOROOT)/src/Make.$(GOARCH)

binaries: 
packages:  pkg/goopt.a

ifndef GOBIN
GOBIN=$(HOME)/bin
endif

# ugly hack to deal with whitespaces in $GOBIN
nullstring :=
space := $(nullstring) # a space at the end
bindir=$(subst $(space),\ ,$(GOBIN))
pkgdir=$(subst $(space),\ ,$(GOROOT)/pkg/$(GOOS)_$(GOARCH))

.PHONY: test binaries packages install installbins installpkgs $(EXTRAPHONY)
.SUFFIXES: .$(O) .go .got .gotgo $(EXTRASUFFIXES)

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`
.got.gotgo:
	gotgo "$<"

ifneq ($(strip $(shell which gotgo)),)
# looks like we require pkg/gotgo/slice.got as installed package...
pkg/gotgo/slice(string).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p pkg/gotgo/
	$< 'string' > "$@"
endif
pkg/goopt.$(O): pkg/goopt.go pkg/gotgo/slice(string).$(O)
pkg/goopt.a: pkg/goopt.$(O)
	gopack grc $@ $<
$(pkgdir)/goopt.a: pkg/goopt.a
	mkdir -p $(pkgdir)/
	cp $< $@


pkg/gotgo/slice(string).$(O): pkg/gotgo/slice(string).go

testit: testit.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
testit.$(O): testit.go

installbins: 
installpkgs:  $(pkgdir)/goopt.a
