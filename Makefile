# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

all: Makefile packages

Makefile: scripts/make.header $(wildcard *.go)
	cp -f scripts/make.header $@
	gotmake >> $@

test: install testit

install: installbins installpkgs


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

.PHONY: test binaries packages install installbins installpkgs
.SUFFIXES: .$(O) .go .got .gotgo

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`
.got.gotgo:
	gotgo "$<"

# looks like we require pkg/gotgo/slice.got as installed package...
pkg/gotgo/slice(string).go: $(pkgdir)/./gotgo/slice.gotgo
	mkdir -p pkg/gotgo/
	$< 'string' > "$@"
pkg/goopt.$(O): pkg/goopt.go pkg/gotgo/slice(string).$(O)
pkg/goopt.a: pkg/goopt.$(O)
	gopack grc $@ $<
$(pkgdir)/goopt.a: pkg/goopt.a
	mkdir -p $(pkgdir)/
	cp $< $@


testit: testit.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
testit.$(O): testit.go

installbins: 
installpkgs:  $(pkgdir)/goopt.a
