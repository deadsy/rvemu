
DIRS = softfp cmd test

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done
	-rm .rv32emu_history
	-rm .rv64emu_history
	-rm log.txt
