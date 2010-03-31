# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.$(GOARCH)

ifndef GOBIN
GOBIN=$(HOME)/bin
endif

# ugly hack to deal with whitespaces in $GOBIN
nullstring :=
space := $(nullstring) # a space at the end
bindir=$(subst $(space),\ ,$(GOBIN))
pkgdir=$(subst $(space),\ ,$(GOROOT)/pkg/$(GOOS)_$(GOARCH))
srcpkgdir=$(subst $(space),\ ,$(GOROOT)/src/pkg)

.PHONY: test install clean
.SUFFIXES: .$(O) .go

.go.$(O):
	cd `dirname "$<"`; $(GC) `basename "$<"`

all: pkg/goopt.a
test: install testit
install: $(pkgdir)/goopt.a
clean:
	rm -f *.$(O) */*.$(O) pkg/*.a testit

ifneq ($(strip $(shell which gotgo)),)
pkg/slice.go: $(srcpkgdir)/gotgo/slice.got
	gotgo --package-name goopt -o "$@" "$<" string
endif

pkg/goopt.$(O): pkg/goopt.go pkg/slice.go
	$(GC) -o $@ $^ 
pkg/goopt.a: pkg/goopt.$(O)
	gopack grc $@ $<
$(pkgdir)/goopt.a: pkg/goopt.a
	mkdir -p $(pkgdir)/
	cp $< $@

testit: testit.$(O)
	@mkdir -p bin
	$(LD) -o $@ $<
testit.$(O): testit.go

