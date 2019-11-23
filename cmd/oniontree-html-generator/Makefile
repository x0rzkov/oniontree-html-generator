TARGET=oniontree-html-generator
LDFLAGS=-w -s

.PHONY: all rebuild clean

all:
	go build -v -ldflags="$(LDFLAGS)" -o "$(TARGET)"

rebuild:
	go build -a -v -ldflags="$(LDFLAGS)" -o "$(TARGET)"

clean:
	$(RM) $(TARGET)

