
DIRS = softfp cmd test

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done
	-rm log.txt
