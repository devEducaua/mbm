PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

all: build

build:
	go build .

install:
	mv $(TARGET) $(BINDIR)

uninstall:
	rm $(BINDIR)/$(TARGET)
