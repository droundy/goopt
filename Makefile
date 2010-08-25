# Copyright 2010 David Roundy, roundyd@physics.oregonstate.edu.
# All rights reserved.

include $(GOROOT)/src/Make.inc

TARG=github.com/droundy/goopt

GOFILES=\
        goopt.go\
        slice.go\

include $(GOROOT)/src/Make.pkg

# ifneq ($(strip $(shell which gotgo)),)
# pkg/slice.go: $(srcpkgdir)/gotgo/slice.got
# 	gotgo --package-name goopt -o "$@" "$<" string
# endif
