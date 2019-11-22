
OUTPUT = test2

XTOOLS = /opt/riscv64-unknown-elf-gcc-8.3.0-2019.08.0-x86_64-linux-ubuntu14/bin/riscv64-unknown-elf-

X_GCC = $(XTOOLS)gcc

SRC = main.c

OBJ = $(patsubst %.c, %.o, $(SRC))

# include paths
INCLUDE =

# defines
DEFINE =

# compiler flags
#X_CFLAGS = -march=rv32g -mabi=ilp32d
#X_CFLAGS = -march=rv64g -mabi=lp64d
X_CFLAGS = -march=rv32gc -mabi=ilp32d
#X_CFLAGS = -march=rv64gc -mabi=lp64d
X_CFLAGS += -Wall -Wstrict-prototypes -Werror

# linker flags
X_LDFLAGS = -Wl,-Map,$(OUTPUT).map -Wl,--gc-sections

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
