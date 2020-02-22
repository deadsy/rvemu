
DIRS = compliance \
       decode \
       rvemu \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done
