
DIRS = gen \
       da \
       rv32emu \
       rv64emu \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done
