
DIRS = cmd \
       test \

all:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

clean:
	for dir in $(DIRS); do \
		$(MAKE) -C ./$$dir $@; \
	done

datest:
	./cmd/da/da -f ./test/dump/rv32imc.txt
