
OUTPUT = test1

XTOOLS = /opt/riscv64-unknown-elf-gcc-8.3.0-2019.08.0-x86_64-linux-ubuntu14/bin/riscv64-unknown-elf-

X_GCC = $(XTOOLS)gcc
X_OBJCOPY = $(XTOOLS)objcopy

SRC = main.c

OBJ = $(patsubst %.c, %.o, $(SRC))
OBJ += start.o

# include paths
INCLUDE =

# defines
DEFINE =

# compiler flags
X_CFLAGS = -march=rv32gc -mabi=ilp32
X_CFLAGS += -Wall -Wstrict-prototypes -Werror

# linker flags
LDSCRIPT = emu.ld
X_LDFLAGS = -T$(LDSCRIPT) -Wl,-Map,$(OUTPUT).map -Wl,--gc-sections -nostdlib

.S.o:
	$(X_GCC) $(INCLUDE) $(DEFINE) $(X_CFLAGS) -c $< -o $@
.c.o:
	$(X_GCC) $(INCLUDE) $(DEFINE) $(X_CFLAGS) -c $< -o $@

.PHONY: all clean

all: $(OBJ)
	$(X_GCC) $(X_CFLAGS) $(X_LDFLAGS) $(OBJ) -o $(OUTPUT)

clean:
	-rm $(OBJ)	
	-rm $(OUTPUT)
	-rm $(OUTPUT).map	
	-rm $(OUTPUT).bin	
