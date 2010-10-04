include $(GOROOT)/src/Make.inc

TARG=goopt-example
GOFILES=example.go

everything: goopt-example.html goopt-example

include $(GOROOT)/src/Make.cmd

# The following install works for Debian (and probably other linux
# distributions), but you probably want to be more sophisticated for
# your real program.
/etc/bash_completion.d/goopt-example: example.bash.completion
	cp -f $< $@

man: man/man1/goopt-example.1
installman: man/man1/goopt-example.1
	cp -f $? /usr/share/man/man1/

# The following rule is generic.
man/man1/%.1: %
	@mkdir -p `dirname $@`
	./$< --create-manpage > $@

%.html: man/man1/%.1
	cat header.html > $@
	groff -man -Thtml $< | tail -n +19 >> $@
