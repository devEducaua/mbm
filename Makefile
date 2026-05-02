PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

TARGET = mbm
TARGETDIR = bin
MANPAGE = mbm.1
SRCS = $(wildcard *.go)

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$(TARGET) .

install:
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/
	cp ./$(MANPAGE) $(MANDIR)/
	gzip -f $(MANDIR)/$(MANPAGE)

uninstall:
	rm $(BINDIR)/$(TARGET)
	rm -f $(MANDIR)/$(MANPAGE).gz

clean: 
	rm -r $(TARGETDIR)

