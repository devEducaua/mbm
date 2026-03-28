PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
TARGET = mbm

all: build

build:
	go build .

install:
	mv $(TARGET) $(BINDIR)

uninstall:
	rm $(BINDIR)/$(TARGET)
