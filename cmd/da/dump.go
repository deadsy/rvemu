//-----------------------------------------------------------------------------
/*

RISC-V Disassembler Test Dumps

Derived from running objdump on gcc *.o files.

*/
//-----------------------------------------------------------------------------

package main

//-----------------------------------------------------------------------------

const dump_rv32im = `

00000000 <nybble>:
   0:   fe010113                addi    sp,sp,-32
   4:   00812e23                sw      s0,28(sp)
   8:   02010413                addi    s0,sp,32
   c:   00050793                mv      a5,a0
  10:   fef407a3                sb      a5,-17(s0)
  14:   fef44783                lbu     a5,-17(s0)
  18:   00f7f793                andi    a5,a5,15
  1c:   fef407a3                sb      a5,-17(s0)
  20:   fef44703                lbu     a4,-17(s0)
  24:   00900793                li      a5,9
  28:   02e7f063                bgeu    a5,a4,48 <.L2>
  2c:   fef44703                lbu     a4,-17(s0)
  30:   00f00793                li      a5,15
  34:   00e7ea63                bltu    a5,a4,48 <.L2>
  38:   fef44783                lbu     a5,-17(s0)
  3c:   05778793                addi    a5,a5,87
  40:   0ff7f793                andi    a5,a5,255
  44:   0100006f                j       54 <.L3>

00000048 <.L2>:
  48:   fef44783                lbu     a5,-17(s0)
  4c:   03078793                addi    a5,a5,48
  50:   0ff7f793                andi    a5,a5,255

00000054 <.L3>:
  54:   00078513                mv      a0,a5
  58:   01c12403                lw      s0,28(sp)
  5c:   02010113                addi    sp,sp,32
  60:   00008067                ret

00000064 <hex8>:
  64:   fe010113                addi    sp,sp,-32
  68:   00112e23                sw      ra,28(sp)
  6c:   00812c23                sw      s0,24(sp)
  70:   00912a23                sw      s1,20(sp)
  74:   02010413                addi    s0,sp,32
  78:   fea42623                sw      a0,-20(s0)
  7c:   00058793                mv      a5,a1
  80:   fef405a3                sb      a5,-21(s0)
  84:   feb44783                lbu     a5,-21(s0)
  88:   0047d793                srli    a5,a5,0x4
  8c:   0ff7f793                andi    a5,a5,255
  90:   00078513                mv      a0,a5
  94:   00000097                auipc   ra,0x0
  98:   000080e7                jalr    ra # 94 <hex8+0x30>
  9c:   00050793                mv      a5,a0
  a0:   00078713                mv      a4,a5
  a4:   fec42783                lw      a5,-20(s0)
  a8:   00e78023                sb      a4,0(a5)
  ac:   fec42783                lw      a5,-20(s0)
  b0:   00178493                addi    s1,a5,1
  b4:   feb44783                lbu     a5,-21(s0)
  b8:   00078513                mv      a0,a5
  bc:   00000097                auipc   ra,0x0
  c0:   000080e7                jalr    ra # bc <hex8+0x58>
  c4:   00050793                mv      a5,a0
  c8:   00f48023                sb      a5,0(s1)
  cc:   fec42783                lw      a5,-20(s0)
  d0:   00278793                addi    a5,a5,2
  d4:   00078023                sb      zero,0(a5)
  d8:   fec42783                lw      a5,-20(s0)
  dc:   00078513                mv      a0,a5
  e0:   01c12083                lw      ra,28(sp)
  e4:   01812403                lw      s0,24(sp)
  e8:   01412483                lw      s1,20(sp)
  ec:   02010113                addi    sp,sp,32
  f0:   00008067                ret

000000f4 <hex16>:
  f4:   fe010113                addi    sp,sp,-32
  f8:   00112e23                sw      ra,28(sp)
  fc:   00812c23                sw      s0,24(sp)
 100:   02010413                addi    s0,sp,32
 104:   fea42623                sw      a0,-20(s0)
 108:   00058793                mv      a5,a1
 10c:   fef41523                sh      a5,-22(s0)
 110:   fea45783                lhu     a5,-22(s0)
 114:   0087d793                srli    a5,a5,0x8
 118:   01079793                slli    a5,a5,0x10
 11c:   0107d793                srli    a5,a5,0x10
 120:   0ff7f793                andi    a5,a5,255
 124:   00078593                mv      a1,a5
 128:   fec42503                lw      a0,-20(s0)
 12c:   00000097                auipc   ra,0x0
 130:   000080e7                jalr    ra # 12c <hex16+0x38>
 134:   fec42783                lw      a5,-20(s0)
 138:   00278793                addi    a5,a5,2
 13c:   fea45703                lhu     a4,-22(s0)
 140:   0ff77713                andi    a4,a4,255
 144:   00070593                mv      a1,a4
 148:   00078513                mv      a0,a5
 14c:   00000097                auipc   ra,0x0
 150:   000080e7                jalr    ra # 14c <hex16+0x58>
 154:   fec42783                lw      a5,-20(s0)
 158:   00078513                mv      a0,a5
 15c:   01c12083                lw      ra,28(sp)
 160:   01812403                lw      s0,24(sp)
 164:   02010113                addi    sp,sp,32
 168:   00008067                ret

0000016c <hex32>:
 16c:   fe010113                addi    sp,sp,-32
 170:   00112e23                sw      ra,28(sp)
 174:   00812c23                sw      s0,24(sp)
 178:   02010413                addi    s0,sp,32
 17c:   fea42623                sw      a0,-20(s0)
 180:   feb42423                sw      a1,-24(s0)
 184:   fe842783                lw      a5,-24(s0)
 188:   0107d793                srli    a5,a5,0x10
 18c:   01079793                slli    a5,a5,0x10
 190:   0107d793                srli    a5,a5,0x10
 194:   00078593                mv      a1,a5
 198:   fec42503                lw      a0,-20(s0)
 19c:   00000097                auipc   ra,0x0
 1a0:   000080e7                jalr    ra # 19c <hex32+0x30>
 1a4:   fec42783                lw      a5,-20(s0)
 1a8:   00478793                addi    a5,a5,4
 1ac:   fe842703                lw      a4,-24(s0)
 1b0:   01071713                slli    a4,a4,0x10
 1b4:   01075713                srli    a4,a4,0x10
 1b8:   00070593                mv      a1,a4
 1bc:   00078513                mv      a0,a5
 1c0:   00000097                auipc   ra,0x0
 1c4:   000080e7                jalr    ra # 1c0 <hex32+0x54>
 1c8:   fec42783                lw      a5,-20(s0)
 1cc:   00078513                mv      a0,a5
 1d0:   01c12083                lw      ra,28(sp)
 1d4:   01812403                lw      s0,24(sp)
 1d8:   02010113                addi    sp,sp,32
 1dc:   00008067                ret

000001e0 <itoa>:
 1e0:   fd010113                addi    sp,sp,-48
 1e4:   02812623                sw      s0,44(sp)
 1e8:   03010413                addi    s0,sp,48
 1ec:   fca42e23                sw      a0,-36(s0)
 1f0:   fcb42c23                sw      a1,-40(s0)
 1f4:   fe042423                sw      zero,-24(s0)
 1f8:   fe042223                sw      zero,-28(s0)
 1fc:   fd842783                lw      a5,-40(s0)
 200:   0007da63                bgez    a5,214 <.L11>
 204:   fd842783                lw      a5,-40(s0)
 208:   40f007b3                neg     a5,a5
 20c:   fef42623                sw      a5,-20(s0)
 210:   00c0006f                j       21c <.L13>

00000214 <.L11>:
 214:   fd842783                lw      a5,-40(s0)
 218:   fef42623                sw      a5,-20(s0)

0000021c <.L13>:
 21c:   fec42703                lw      a4,-20(s0)
 220:   00a00793                li      a5,10
 224:   02f777b3                remu    a5,a4,a5
 228:   0ff7f713                andi    a4,a5,255
 22c:   fe842783                lw      a5,-24(s0)
 230:   00178693                addi    a3,a5,1
 234:   fed42423                sw      a3,-24(s0)
 238:   00078693                mv      a3,a5
 23c:   fdc42783                lw      a5,-36(s0)
 240:   00d787b3                add     a5,a5,a3
 244:   03070713                addi    a4,a4,48
 248:   0ff77713                andi    a4,a4,255
 24c:   00e78023                sb      a4,0(a5)
 250:   fec42703                lw      a4,-20(s0)
 254:   00a00793                li      a5,10
 258:   02f757b3                divu    a5,a4,a5
 25c:   fef42623                sw      a5,-20(s0)
 260:   fec42783                lw      a5,-20(s0)
 264:   fa079ce3                bnez    a5,21c <.L13>
 268:   fd842783                lw      a5,-40(s0)
 26c:   0207d263                bgez    a5,290 <.L14>
 270:   fe842783                lw      a5,-24(s0)
 274:   00178713                addi    a4,a5,1
 278:   fee42423                sw      a4,-24(s0)
 27c:   00078713                mv      a4,a5
 280:   fdc42783                lw      a5,-36(s0)
 284:   00e787b3                add     a5,a5,a4
 288:   02d00713                li      a4,45
 28c:   00e78023                sb      a4,0(a5)

00000290 <.L14>:
 290:   fe842783                lw      a5,-24(s0)
 294:   fdc42703                lw      a4,-36(s0)
 298:   00f707b3                add     a5,a4,a5
 29c:   00078023                sb      zero,0(a5)
 2a0:   fe842783                lw      a5,-24(s0)
 2a4:   fff78793                addi    a5,a5,-1
 2a8:   fef42423                sw      a5,-24(s0)
 2ac:   0640006f                j       310 <.L15>

000002b0 <.L16>:
 2b0:   fe442783                lw      a5,-28(s0)
 2b4:   fdc42703                lw      a4,-36(s0)
 2b8:   00f707b3                add     a5,a4,a5
 2bc:   0007c783                lbu     a5,0(a5)
 2c0:   fef401a3                sb      a5,-29(s0)
 2c4:   fe842783                lw      a5,-24(s0)
 2c8:   fdc42703                lw      a4,-36(s0)
 2cc:   00f70733                add     a4,a4,a5
 2d0:   fe442783                lw      a5,-28(s0)
 2d4:   00178693                addi    a3,a5,1
 2d8:   fed42223                sw      a3,-28(s0)
 2dc:   00078693                mv      a3,a5
 2e0:   fdc42783                lw      a5,-36(s0)
 2e4:   00d787b3                add     a5,a5,a3
 2e8:   00074703                lbu     a4,0(a4)
 2ec:   00e78023                sb      a4,0(a5)
 2f0:   fe842783                lw      a5,-24(s0)
 2f4:   fff78713                addi    a4,a5,-1
 2f8:   fee42423                sw      a4,-24(s0)
 2fc:   00078713                mv      a4,a5
 300:   fdc42783                lw      a5,-36(s0)
 304:   00e787b3                add     a5,a5,a4
 308:   fe344703                lbu     a4,-29(s0)
 30c:   00e78023                sb      a4,0(a5)

00000310 <.L15>:
 310:   fe442703                lw      a4,-28(s0)
 314:   fe842783                lw      a5,-24(s0)
 318:   f8f74ce3                blt     a4,a5,2b0 <.L16>
 31c:   fdc42783                lw      a5,-36(s0)
 320:   00078513                mv      a0,a5
 324:   02c12403                lw      s0,44(sp)
 328:   03010113                addi    sp,sp,48
 32c:   00008067                ret

00000330 <main>:
 330:   fd010113                addi    sp,sp,-48
 334:   02112623                sw      ra,44(sp)
 338:   02812423                sw      s0,40(sp)
 33c:   03010413                addi    s0,sp,48
 340:   fd040793                addi    a5,s0,-48
 344:   00000593                li      a1,0
 348:   00078513                mv      a0,a5
 34c:   00000097                auipc   ra,0x0
 350:   000080e7                jalr    ra # 34c <main+0x1c>
 354:   00050793                mv      a5,a0
 358:   00078513                mv      a0,a5
 35c:   00000097                auipc   ra,0x0
 360:   000080e7                jalr    ra # 35c <main+0x2c>
 364:   fd040793                addi    a5,s0,-48
 368:   4d200593                li      a1,1234
 36c:   00078513                mv      a0,a5
 370:   00000097                auipc   ra,0x0
 374:   000080e7                jalr    ra # 370 <main+0x40>
 378:   00050793                mv      a5,a0
 37c:   00078513                mv      a0,a5
 380:   00000097                auipc   ra,0x0
 384:   000080e7                jalr    ra # 380 <main+0x50>
 388:   fd040793                addi    a5,s0,-48
 38c:   b2e00593                li      a1,-1234
 390:   00078513                mv      a0,a5
 394:   00000097                auipc   ra,0x0
 398:   000080e7                jalr    ra # 394 <main+0x64>
 39c:   00050793                mv      a5,a0
 3a0:   00078513                mv      a0,a5
 3a4:   00000097                auipc   ra,0x0
 3a8:   000080e7                jalr    ra # 3a4 <main+0x74>
 3ac:   fd040713                addi    a4,s0,-48
 3b0:   800007b7                lui     a5,0x80000
 3b4:   fff7c593                not     a1,a5
 3b8:   00070513                mv      a0,a4
 3bc:   00000097                auipc   ra,0x0
 3c0:   000080e7                jalr    ra # 3bc <main+0x8c>
 3c4:   00050793                mv      a5,a0
 3c8:   00078513                mv      a0,a5
 3cc:   00000097                auipc   ra,0x0
 3d0:   000080e7                jalr    ra # 3cc <main+0x9c>
 3d4:   fd040793                addi    a5,s0,-48
 3d8:   800005b7                lui     a1,0x80000
 3dc:   00078513                mv      a0,a5
 3e0:   00000097                auipc   ra,0x0
 3e4:   000080e7                jalr    ra # 3e0 <main+0xb0>
 3e8:   00050793                mv      a5,a0
 3ec:   00078513                mv      a0,a5
 3f0:   00000097                auipc   ra,0x0
 3f4:   000080e7                jalr    ra # 3f0 <main+0xc0>
 3f8:   fd040793                addi    a5,s0,-48
 3fc:   0ab00593                li      a1,171
 400:   00078513                mv      a0,a5
 404:   00000097                auipc   ra,0x0
 408:   000080e7                jalr    ra # 404 <main+0xd4>
 40c:   00050793                mv      a5,a0
 410:   00078513                mv      a0,a5
 414:   00000097                auipc   ra,0x0
 418:   000080e7                jalr    ra # 414 <main+0xe4>
 41c:   fd040713                addi    a4,s0,-48
 420:   0000b7b7                lui     a5,0xb
 424:   bcd78593                addi    a1,a5,-1075 # abcd <main+0xa89d>
 428:   00070513                mv      a0,a4
 42c:   00000097                auipc   ra,0x0
 430:   000080e7                jalr    ra # 42c <main+0xfc>
 434:   00050793                mv      a5,a0
 438:   00078513                mv      a0,a5
 43c:   00000097                auipc   ra,0x0
 440:   000080e7                jalr    ra # 43c <main+0x10c>
 444:   fd040713                addi    a4,s0,-48
 448:   deadc7b7                lui     a5,0xdeadc
 44c:   eef78593                addi    a1,a5,-273 # deadbeef <main+0xdeadbbbf>
 450:   00070513                mv      a0,a4
 454:   00000097                auipc   ra,0x0
 458:   000080e7                jalr    ra # 454 <main+0x124>
 45c:   00050793                mv      a5,a0
 460:   00078513                mv      a0,a5
 464:   00000097                auipc   ra,0x0
 468:   000080e7                jalr    ra # 464 <main+0x134>
 46c:   00000793                li      a5,0
 470:   00078513                mv      a0,a5
 474:   02c12083                lw      ra,44(sp)
 478:   02812403                lw      s0,40(sp)
 47c:   03010113                addi    sp,sp,48
 480:   00008067                ret
`

//-----------------------------------------------------------------------------

const dump1_rv32imc = `

00000000 <nybble>:
   0:   1101                    addi    sp,sp,-32
   2:   ce22                    sw      s0,28(sp)
   4:   1000                    addi    s0,sp,32
   6:   87aa                    mv      a5,a0
   8:   fef407a3                sb      a5,-17(s0)
   c:   fef44783                lbu     a5,-17(s0)
  10:   8bbd                    andi    a5,a5,15
  12:   fef407a3                sb      a5,-17(s0)
  16:   fef44703                lbu     a4,-17(s0)
  1a:   47a5                    li      a5,9
  1c:   00e7fe63                bgeu    a5,a4,38 <.L2>
  20:   fef44703                lbu     a4,-17(s0)
  24:   47bd                    li      a5,15
  26:   00e7e963                bltu    a5,a4,38 <.L2>
  2a:   fef44783                lbu     a5,-17(s0)
  2e:   05778793                addi    a5,a5,87
  32:   0ff7f793                andi    a5,a5,255
  36:   a039                    j       44 <.L3>

00000038 <.L2>:
  38:   fef44783                lbu     a5,-17(s0)
  3c:   03078793                addi    a5,a5,48
  40:   0ff7f793                andi    a5,a5,255

00000044 <.L3>:
  44:   853e                    mv      a0,a5
  46:   4472                    lw      s0,28(sp)
  48:   6105                    addi    sp,sp,32
  4a:   8082                    ret

0000004c <hex8>:
  4c:   1101                    addi    sp,sp,-32
  4e:   ce06                    sw      ra,28(sp)
  50:   cc22                    sw      s0,24(sp)
  52:   ca26                    sw      s1,20(sp)
  54:   1000                    addi    s0,sp,32
  56:   fea42623                sw      a0,-20(s0)
  5a:   87ae                    mv      a5,a1
  5c:   fef405a3                sb      a5,-21(s0)
  60:   feb44783                lbu     a5,-21(s0)
  64:   8391                    srli    a5,a5,0x4
  66:   0ff7f793                andi    a5,a5,255
  6a:   853e                    mv      a0,a5
  6c:   00000097                auipc   ra,0x0
  70:   000080e7                jalr    ra # 6c <hex8+0x20>
  74:   87aa                    mv      a5,a0
  76:   873e                    mv      a4,a5
  78:   fec42783                lw      a5,-20(s0)
  7c:   00e78023                sb      a4,0(a5)
  80:   fec42783                lw      a5,-20(s0)
  84:   00178493                addi    s1,a5,1
  88:   feb44783                lbu     a5,-21(s0)
  8c:   853e                    mv      a0,a5
  8e:   00000097                auipc   ra,0x0
  92:   000080e7                jalr    ra # 8e <hex8+0x42>
  96:   87aa                    mv      a5,a0
  98:   00f48023                sb      a5,0(s1)
  9c:   fec42783                lw      a5,-20(s0)
  a0:   0789                    addi    a5,a5,2
  a2:   00078023                sb      zero,0(a5)
  a6:   fec42783                lw      a5,-20(s0)
  aa:   853e                    mv      a0,a5
  ac:   40f2                    lw      ra,28(sp)
  ae:   4462                    lw      s0,24(sp)
  b0:   44d2                    lw      s1,20(sp)
  b2:   6105                    addi    sp,sp,32
  b4:   8082                    ret

000000b6 <hex16>:
  b6:   1101                    addi    sp,sp,-32
  b8:   ce06                    sw      ra,28(sp)
  ba:   cc22                    sw      s0,24(sp)
  bc:   1000                    addi    s0,sp,32
  be:   fea42623                sw      a0,-20(s0)
  c2:   87ae                    mv      a5,a1
  c4:   fef41523                sh      a5,-22(s0)
  c8:   fea45783                lhu     a5,-22(s0)
  cc:   83a1                    srli    a5,a5,0x8
  ce:   07c2                    slli    a5,a5,0x10
  d0:   83c1                    srli    a5,a5,0x10
  d2:   0ff7f793                andi    a5,a5,255
  d6:   85be                    mv      a1,a5
  d8:   fec42503                lw      a0,-20(s0)
  dc:   00000097                auipc   ra,0x0
  e0:   000080e7                jalr    ra # dc <hex16+0x26>
  e4:   fec42783                lw      a5,-20(s0)
  e8:   0789                    addi    a5,a5,2
  ea:   fea45703                lhu     a4,-22(s0)
  ee:   0ff77713                andi    a4,a4,255
  f2:   85ba                    mv      a1,a4
  f4:   853e                    mv      a0,a5
  f6:   00000097                auipc   ra,0x0
  fa:   000080e7                jalr    ra # f6 <hex16+0x40>
  fe:   fec42783                lw      a5,-20(s0)
 102:   853e                    mv      a0,a5
 104:   40f2                    lw      ra,28(sp)
 106:   4462                    lw      s0,24(sp)
 108:   6105                    addi    sp,sp,32
 10a:   8082                    ret

0000010c <hex32>:
 10c:   1101                    addi    sp,sp,-32
 10e:   ce06                    sw      ra,28(sp)
 110:   cc22                    sw      s0,24(sp)
 112:   1000                    addi    s0,sp,32
 114:   fea42623                sw      a0,-20(s0)
 118:   feb42423                sw      a1,-24(s0)
 11c:   fe842783                lw      a5,-24(s0)
 120:   83c1                    srli    a5,a5,0x10
 122:   07c2                    slli    a5,a5,0x10
 124:   83c1                    srli    a5,a5,0x10
 126:   85be                    mv      a1,a5
 128:   fec42503                lw      a0,-20(s0)
 12c:   00000097                auipc   ra,0x0
 130:   000080e7                jalr    ra # 12c <hex32+0x20>
 134:   fec42783                lw      a5,-20(s0)
 138:   0791                    addi    a5,a5,4
 13a:   fe842703                lw      a4,-24(s0)
 13e:   0742                    slli    a4,a4,0x10
 140:   8341                    srli    a4,a4,0x10
 142:   85ba                    mv      a1,a4
 144:   853e                    mv      a0,a5
 146:   00000097                auipc   ra,0x0
 14a:   000080e7                jalr    ra # 146 <hex32+0x3a>
 14e:   fec42783                lw      a5,-20(s0)
 152:   853e                    mv      a0,a5
 154:   40f2                    lw      ra,28(sp)
 156:   4462                    lw      s0,24(sp)
 158:   6105                    addi    sp,sp,32
 15a:   8082                    ret

0000015c <itoa>:
 15c:   7179                    addi    sp,sp,-48
 15e:   d622                    sw      s0,44(sp)
 160:   1800                    addi    s0,sp,48
 162:   fca42e23                sw      a0,-36(s0)
 166:   fcb42c23                sw      a1,-40(s0)
 16a:   fe042423                sw      zero,-24(s0)
 16e:   fe042223                sw      zero,-28(s0)
 172:   fd842783                lw      a5,-40(s0)
 176:   0007d963                bgez    a5,188 <.L11>
 17a:   fd842783                lw      a5,-40(s0)
 17e:   40f007b3                neg     a5,a5
 182:   fef42623                sw      a5,-20(s0)
 186:   a029                    j       190 <.L13>

00000188 <.L11>:
 188:   fd842783                lw      a5,-40(s0)
 18c:   fef42623                sw      a5,-20(s0)

00000190 <.L13>:
 190:   fec42703                lw      a4,-20(s0)
 194:   47a9                    li      a5,10
 196:   02f777b3                remu    a5,a4,a5
 19a:   0ff7f713                andi    a4,a5,255
 19e:   fe842783                lw      a5,-24(s0)
 1a2:   00178693                addi    a3,a5,1
 1a6:   fed42423                sw      a3,-24(s0)
 1aa:   86be                    mv      a3,a5
 1ac:   fdc42783                lw      a5,-36(s0)
 1b0:   97b6                    add     a5,a5,a3
 1b2:   03070713                addi    a4,a4,48
 1b6:   0ff77713                andi    a4,a4,255
 1ba:   00e78023                sb      a4,0(a5)
 1be:   fec42703                lw      a4,-20(s0)
 1c2:   47a9                    li      a5,10
 1c4:   02f757b3                divu    a5,a4,a5
 1c8:   fef42623                sw      a5,-20(s0)
 1cc:   fec42783                lw      a5,-20(s0)
 1d0:   f3e1                    bnez    a5,190 <.L13>
 1d2:   fd842783                lw      a5,-40(s0)
 1d6:   0207d063                bgez    a5,1f6 <.L14>
 1da:   fe842783                lw      a5,-24(s0)
 1de:   00178713                addi    a4,a5,1
 1e2:   fee42423                sw      a4,-24(s0)
 1e6:   873e                    mv      a4,a5
 1e8:   fdc42783                lw      a5,-36(s0)
 1ec:   97ba                    add     a5,a5,a4
 1ee:   02d00713                li      a4,45
 1f2:   00e78023                sb      a4,0(a5)

000001f6 <.L14>:
 1f6:   fe842783                lw      a5,-24(s0)
 1fa:   fdc42703                lw      a4,-36(s0)
 1fe:   97ba                    add     a5,a5,a4
 200:   00078023                sb      zero,0(a5)
 204:   fe842783                lw      a5,-24(s0)
 208:   17fd                    addi    a5,a5,-1
 20a:   fef42423                sw      a5,-24(s0)
 20e:   a899                    j       264 <.L15>

00000210 <.L16>:
 210:   fe442783                lw      a5,-28(s0)
 214:   fdc42703                lw      a4,-36(s0)
 218:   97ba                    add     a5,a5,a4
 21a:   0007c783                lbu     a5,0(a5)
 21e:   fef401a3                sb      a5,-29(s0)
 222:   fe842783                lw      a5,-24(s0)
 226:   fdc42703                lw      a4,-36(s0)
 22a:   973e                    add     a4,a4,a5
 22c:   fe442783                lw      a5,-28(s0)
 230:   00178693                addi    a3,a5,1
 234:   fed42223                sw      a3,-28(s0)
 238:   86be                    mv      a3,a5
 23a:   fdc42783                lw      a5,-36(s0)
 23e:   97b6                    add     a5,a5,a3
 240:   00074703                lbu     a4,0(a4)
 244:   00e78023                sb      a4,0(a5)
 248:   fe842783                lw      a5,-24(s0)
 24c:   fff78713                addi    a4,a5,-1
 250:   fee42423                sw      a4,-24(s0)
 254:   873e                    mv      a4,a5
 256:   fdc42783                lw      a5,-36(s0)
 25a:   97ba                    add     a5,a5,a4
 25c:   fe344703                lbu     a4,-29(s0)
 260:   00e78023                sb      a4,0(a5)

00000264 <.L15>:
 264:   fe442703                lw      a4,-28(s0)
 268:   fe842783                lw      a5,-24(s0)
 26c:   faf742e3                blt     a4,a5,210 <.L16>
 270:   fdc42783                lw      a5,-36(s0)
 274:   853e                    mv      a0,a5
 276:   5432                    lw      s0,44(sp)
 278:   6145                    addi    sp,sp,48
 27a:   8082                    ret

0000027c <main>:
 27c:   7179                    addi    sp,sp,-48
 27e:   d606                    sw      ra,44(sp)
 280:   d422                    sw      s0,40(sp)
 282:   1800                    addi    s0,sp,48
 284:   fd040793                addi    a5,s0,-48
 288:   4581                    li      a1,0
 28a:   853e                    mv      a0,a5
 28c:   00000097                auipc   ra,0x0
 290:   000080e7                jalr    ra # 28c <main+0x10>
 294:   87aa                    mv      a5,a0
 296:   853e                    mv      a0,a5
 298:   00000097                auipc   ra,0x0
 29c:   000080e7                jalr    ra # 298 <main+0x1c>
 2a0:   fd040793                addi    a5,s0,-48
 2a4:   4d200593                li      a1,1234
 2a8:   853e                    mv      a0,a5
 2aa:   00000097                auipc   ra,0x0
 2ae:   000080e7                jalr    ra # 2aa <main+0x2e>
 2b2:   87aa                    mv      a5,a0
 2b4:   853e                    mv      a0,a5
 2b6:   00000097                auipc   ra,0x0
 2ba:   000080e7                jalr    ra # 2b6 <main+0x3a>
 2be:   fd040793                addi    a5,s0,-48
 2c2:   b2e00593                li      a1,-1234
 2c6:   853e                    mv      a0,a5
 2c8:   00000097                auipc   ra,0x0
 2cc:   000080e7                jalr    ra # 2c8 <main+0x4c>
 2d0:   87aa                    mv      a5,a0
 2d2:   853e                    mv      a0,a5
 2d4:   00000097                auipc   ra,0x0
 2d8:   000080e7                jalr    ra # 2d4 <main+0x58>
 2dc:   fd040713                addi    a4,s0,-48
 2e0:   800007b7                lui     a5,0x80000
 2e4:   fff7c593                not     a1,a5
 2e8:   853a                    mv      a0,a4
 2ea:   00000097                auipc   ra,0x0
 2ee:   000080e7                jalr    ra # 2ea <main+0x6e>
 2f2:   87aa                    mv      a5,a0
 2f4:   853e                    mv      a0,a5
 2f6:   00000097                auipc   ra,0x0
 2fa:   000080e7                jalr    ra # 2f6 <main+0x7a>
 2fe:   fd040793                addi    a5,s0,-48
 302:   800005b7                lui     a1,0x80000
 306:   853e                    mv      a0,a5
 308:   00000097                auipc   ra,0x0
 30c:   000080e7                jalr    ra # 308 <main+0x8c>
 310:   87aa                    mv      a5,a0
 312:   853e                    mv      a0,a5
 314:   00000097                auipc   ra,0x0
 318:   000080e7                jalr    ra # 314 <main+0x98>
 31c:   fd040793                addi    a5,s0,-48
 320:   0ab00593                li      a1,171
 324:   853e                    mv      a0,a5
 326:   00000097                auipc   ra,0x0
 32a:   000080e7                jalr    ra # 326 <main+0xaa>
 32e:   87aa                    mv      a5,a0
 330:   853e                    mv      a0,a5
 332:   00000097                auipc   ra,0x0
 336:   000080e7                jalr    ra # 332 <main+0xb6>
 33a:   fd040713                addi    a4,s0,-48
 33e:   67ad                    lui     a5,0xb
 340:   bcd78593                addi    a1,a5,-1075 # abcd <main+0xa951>
 344:   853a                    mv      a0,a4
 346:   00000097                auipc   ra,0x0
 34a:   000080e7                jalr    ra # 346 <main+0xca>
 34e:   87aa                    mv      a5,a0
 350:   853e                    mv      a0,a5
 352:   00000097                auipc   ra,0x0
 356:   000080e7                jalr    ra # 352 <main+0xd6>
 35a:   fd040713                addi    a4,s0,-48
 35e:   deadc7b7                lui     a5,0xdeadc
 362:   eef78593                addi    a1,a5,-273 # deadbeef <main+0xdeadbc73>
 366:   853a                    mv      a0,a4
 368:   00000097                auipc   ra,0x0
 36c:   000080e7                jalr    ra # 368 <main+0xec>
 370:   87aa                    mv      a5,a0
 372:   853e                    mv      a0,a5
 374:   00000097                auipc   ra,0x0
 378:   000080e7                jalr    ra # 374 <main+0xf8>
 37c:   4781                    li      a5,0
 37e:   853e                    mv      a0,a5
 380:   50b2                    lw      ra,44(sp)
 382:   5422                    lw      s0,40(sp)
 384:   6145                    addi    sp,sp,48
 386:   8082                    ret
`

//-----------------------------------------------------------------------------
