include $(GOROOT)/src/Make.$(GOARCH)

TARG=goopt
GOFILES=\
	goopt.go\

include $(GOROOT)/src/Make.pkg
