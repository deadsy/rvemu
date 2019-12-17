#!/bin/bash

SET1="/opt/riscv64-unknown-elf-gcc-8.3.0-2019.08.0-x86_64-linux-ubuntu14/target/share/riscv-tests/isa"
SET2="/home/jasonh/work/riscv/riscv-compliance/work"

echo "****" $SET1 > log.txt
./cmd/compliance/compliance -p $SET1 >> log.txt

echo "****" $SET2 >> log.txt
./cmd/compliance/compliance -p $SET2 >> log.txt
