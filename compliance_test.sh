#!/bin/bash

SET1="/home/jasonh/work/riscv/riscv-tests/isa"
SET2="/home/jasonh/work/riscv/riscv-compliance/work"

echo $SET1 > log.txt
./cmd/compliance/compliance -p $SET1 >> log.txt

echo $SET2 >> log.txt
./cmd/compliance/compliance -p $SET2 >> log.txt
