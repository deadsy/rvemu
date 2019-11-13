
MEMORY
{
  rom : ORIGIN = 0x00000000, LENGTH = 256K
  ram : ORIGIN = 0x80000000, LENGTH = 64K
}

PROVIDE(_ram_start = ORIGIN(ram));
PROVIDE(_ram_end = ORIGIN(ram) + LENGTH(ram));

OUTPUT_ARCH( "riscv" )

ENTRY(_start)

SECTIONS
{
  .text :
  {
    KEEP(*(.text.start))
    *(.text .text.*)
    . = ALIGN(8);
    PROVIDE(_etext = ABSOLUTE(.));
  } >rom

  .rodata :
  {
    *(.rodata .rodata.*)
  } >rom

  .data :
  {
    *(.data .data.*)
    . = ALIGN(8);
    PROVIDE(_edata = ABSOLUTE(.));
  } >ram

  .bss :
  {
    PROVIDE(_bss = ABSOLUTE(.));
    *(.bss .bss.*)
    *(COMMON)
    . = ALIGN(8);
    PROVIDE(_ebss = ABSOLUTE(.));
  } >ram

}
