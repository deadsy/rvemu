# See: https://bellard.org/softfp/
# IEEE754 defines rounding modes and exception flags for FP operations.
# The RISC-V floating point instructions make direct use of these.
# The go runtime doesn't provide easy access to these features, but this
# library does, so we use this library to emulate the FP operations.

VER = 2018-09-15
DIR = softfp-$(VER)
TGZ = $(DIR).tar.gz  
LIB = libsoftfp.a

SRC = $(DIR)/softfloat.c \
			$(DIR)/softfp.c \

OBJ = $(patsubst %.c, %.o, $(SRC))

# include paths
#INCLUDE = -I$(DIR)

# defines 
DEFINE += -DCONFIG_VERSION=\"$(VER)\"

# compiler flags
CFLAGS = -O2 -Wall -Wstrict-prototypes -Werror -MMD

.c.o:
	gcc $(INCLUDE) $(DEFINE) $(CFLAGS) -c $< -o $@

all: untar $(LIB)
	cp $(DIR)/softfp.h .
	cp $(DIR)/cutils.h .

libsoftfp.a: $(OBJ)
	ar rcs $@ $^

.PHONY: untar
untar:
	tar zxf $(TGZ)

.PHONY: clean
clean:
	-rm -rf $(DIR)
	-rm -rf $(LIB)
	-rm -rf *.h
