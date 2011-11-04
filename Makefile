include ${GOROOT}/src/Make.inc

TXT2GO = ./txt2go.sh

TARG = munch
GOFILES = munch.go page.go rss.go rdf.go atom.go httpd.go time.go about.go
CLEANFILES = page.go

include ${GOROOT}/src/Make.cmd

page.go: page.html ${TXT2GO}
	${TXT2GO} page_template_string < page.html > $@
about.go: about.html ${TXT2GO}
	${TXT2GO} about_string < about.html > $@

