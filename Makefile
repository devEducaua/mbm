PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

TARGET = mbm
MANPAGE = mbm.1

all: build

build:
	go build .

install:
	cp $(TARGET) $(BINDIR)/
	cp ./$(MANPAGE) $(MANDIR)/
	gzip -f $(MANDIR)/$(MANPAGE)

uninstall:
	rm $(BINDIR)/$(TARGET)
	rm -f $(MANDIR)/$(MANPAGE).gz

clean: 
	rm $(TARGET)
