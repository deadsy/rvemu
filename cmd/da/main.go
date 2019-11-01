//-----------------------------------------------------------------------------

RISC-V Disassembler

//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/riscv/cpu"
)

//-----------------------------------------------------------------------------

/*

00000000 <nybble>:
   0:	fe010113          	addi	sp,sp,-32
   4:	00812e23          	sw	s0,28(sp)
   8:	02010413          	addi	s0,sp,32
   c:	00050793          	mv	a5,a0
  10:	fef407a3          	sb	a5,-17(s0)
  14:	fef44783          	lbu	a5,-17(s0)
  18:	00f7f793          	andi	a5,a5,15
  1c:	fef407a3          	sb	a5,-17(s0)
  20:	fef44703          	lbu	a4,-17(s0)
  24:	00900793          	li	a5,9
  28:	02e7f063          	bgeu	a5,a4,48 <.L2>
  2c:	fef44703          	lbu	a4,-17(s0)
  30:	00f00793          	li	a5,15
  34:	00e7ea63          	bltu	a5,a4,48 <.L2>
  38:	fef44783          	lbu	a5,-17(s0)
  3c:	05778793          	addi	a5,a5,87
  40:	0ff7f793          	andi	a5,a5,255
  44:	0100006f          	j	54 <.L3>

00000048 <.L2>:
  48:	fef44783          	lbu	a5,-17(s0)
  4c:	03078793          	addi	a5,a5,48
  50:	0ff7f793          	andi	a5,a5,255

00000054 <.L3>:
  54:	00078513          	mv	a0,a5
  58:	01c12403          	lw	s0,28(sp)
  5c:	02010113          	addi	sp,sp,32
  60:	00008067          	ret

00000064 <hex8>:
  64:	fe010113          	addi	sp,sp,-32
  68:	00112e23          	sw	ra,28(sp)
  6c:	00812c23          	sw	s0,24(sp)
  70:	00912a23          	sw	s1,20(sp)
  74:	02010413          	addi	s0,sp,32
  78:	fea42623          	sw	a0,-20(s0)
  7c:	00058793          	mv	a5,a1
  80:	fef405a3          	sb	a5,-21(s0)
  84:	feb44783          	lbu	a5,-21(s0)
  88:	0047d793          	srli	a5,a5,0x4
  8c:	0ff7f793          	andi	a5,a5,255
  90:	00078513          	mv	a0,a5
  94:	00000097          	auipc	ra,0x0
  98:	000080e7          	jalr	ra # 94 <hex8+0x30>
  9c:	00050793          	mv	a5,a0
  a0:	00078713          	mv	a4,a5
  a4:	fec42783          	lw	a5,-20(s0)
  a8:	00e78023          	sb	a4,0(a5)
  ac:	fec42783          	lw	a5,-20(s0)
  b0:	00178493          	addi	s1,a5,1
  b4:	feb44783          	lbu	a5,-21(s0)
  b8:	00078513          	mv	a0,a5
  bc:	00000097          	auipc	ra,0x0
  c0:	000080e7          	jalr	ra # bc <hex8+0x58>
  c4:	00050793          	mv	a5,a0
  c8:	00f48023          	sb	a5,0(s1)
  cc:	fec42783          	lw	a5,-20(s0)
  d0:	00278793          	addi	a5,a5,2
  d4:	00078023          	sb	zero,0(a5)
  d8:	fec42783          	lw	a5,-20(s0)
  dc:	00078513          	mv	a0,a5
  e0:	01c12083          	lw	ra,28(sp)
  e4:	01812403          	lw	s0,24(sp)
  e8:	01412483          	lw	s1,20(sp)
  ec:	02010113          	addi	sp,sp,32
  f0:	00008067          	ret

000000f4 <hex16>:
  f4:	fe010113          	addi	sp,sp,-32
  f8:	00112e23          	sw	ra,28(sp)
  fc:	00812c23          	sw	s0,24(sp)
 100:	02010413          	addi	s0,sp,32
 104:	fea42623          	sw	a0,-20(s0)
 108:	00058793          	mv	a5,a1
 10c:	fef41523          	sh	a5,-22(s0)
 110:	fea45783          	lhu	a5,-22(s0)
 114:	0087d793          	srli	a5,a5,0x8
 118:	01079793          	slli	a5,a5,0x10
 11c:	0107d793          	srli	a5,a5,0x10
 120:	0ff7f793          	andi	a5,a5,255
 124:	00078593          	mv	a1,a5
 128:	fec42503          	lw	a0,-20(s0)
 12c:	00000097          	auipc	ra,0x0
 130:	000080e7          	jalr	ra # 12c <hex16+0x38>
 134:	fec42783          	lw	a5,-20(s0)
 138:	00278793          	addi	a5,a5,2
 13c:	fea45703          	lhu	a4,-22(s0)
 140:	0ff77713          	andi	a4,a4,255
 144:	00070593          	mv	a1,a4
 148:	00078513          	mv	a0,a5
 14c:	00000097          	auipc	ra,0x0
 150:	000080e7          	jalr	ra # 14c <hex16+0x58>
 154:	fec42783          	lw	a5,-20(s0)
 158:	00078513          	mv	a0,a5
 15c:	01c12083          	lw	ra,28(sp)
 160:	01812403          	lw	s0,24(sp)
 164:	02010113          	addi	sp,sp,32
 168:	00008067          	ret

0000016c <hex32>:
 16c:	fe010113          	addi	sp,sp,-32
 170:	00112e23          	sw	ra,28(sp)
 174:	00812c23          	sw	s0,24(sp)
 178:	02010413          	addi	s0,sp,32
 17c:	fea42623          	sw	a0,-20(s0)
 180:	feb42423          	sw	a1,-24(s0)
 184:	fe842783          	lw	a5,-24(s0)
 188:	0107d793          	srli	a5,a5,0x10
 18c:	01079793          	slli	a5,a5,0x10
 190:	0107d793          	srli	a5,a5,0x10
 194:	00078593          	mv	a1,a5
 198:	fec42503          	lw	a0,-20(s0)
 19c:	00000097          	auipc	ra,0x0
 1a0:	000080e7          	jalr	ra # 19c <hex32+0x30>
 1a4:	fec42783          	lw	a5,-20(s0)
 1a8:	00478793          	addi	a5,a5,4
 1ac:	fe842703          	lw	a4,-24(s0)
 1b0:	01071713          	slli	a4,a4,0x10
 1b4:	01075713          	srli	a4,a4,0x10
 1b8:	00070593          	mv	a1,a4
 1bc:	00078513          	mv	a0,a5
 1c0:	00000097          	auipc	ra,0x0
 1c4:	000080e7          	jalr	ra # 1c0 <hex32+0x54>
 1c8:	fec42783          	lw	a5,-20(s0)
 1cc:	00078513          	mv	a0,a5
 1d0:	01c12083          	lw	ra,28(sp)
 1d4:	01812403          	lw	s0,24(sp)
 1d8:	02010113          	addi	sp,sp,32
 1dc:	00008067          	ret

000001e0 <itoa>:
 1e0:	fd010113          	addi	sp,sp,-48
 1e4:	02112623          	sw	ra,44(sp)
 1e8:	02812423          	sw	s0,40(sp)
 1ec:	03010413          	addi	s0,sp,48
 1f0:	fca42e23          	sw	a0,-36(s0)
 1f4:	fcb42c23          	sw	a1,-40(s0)
 1f8:	fe042423          	sw	zero,-24(s0)
 1fc:	fe042223          	sw	zero,-28(s0)
 200:	fd842783          	lw	a5,-40(s0)
 204:	0007da63          	bgez	a5,218 <.L11>
 208:	fd842783          	lw	a5,-40(s0)
 20c:	40f007b3          	neg	a5,a5
 210:	fef42623          	sw	a5,-20(s0)
 214:	00c0006f          	j	220 <.L13>

00000218 <.L11>:
 218:	fd842783          	lw	a5,-40(s0)
 21c:	fef42623          	sw	a5,-20(s0)

00000220 <.L13>:
 220:	fec42783          	lw	a5,-20(s0)
 224:	00a00593          	li	a1,10
 228:	00078513          	mv	a0,a5
 22c:	00000097          	auipc	ra,0x0
 230:	000080e7          	jalr	ra # 22c <.L13+0xc>
 234:	00050793          	mv	a5,a0
 238:	0ff7f713          	andi	a4,a5,255
 23c:	fe842783          	lw	a5,-24(s0)
 240:	00178693          	addi	a3,a5,1
 244:	fed42423          	sw	a3,-24(s0)
 248:	00078693          	mv	a3,a5
 24c:	fdc42783          	lw	a5,-36(s0)
 250:	00d787b3          	add	a5,a5,a3
 254:	03070713          	addi	a4,a4,48
 258:	0ff77713          	andi	a4,a4,255
 25c:	00e78023          	sb	a4,0(a5)
 260:	fec42783          	lw	a5,-20(s0)
 264:	00a00593          	li	a1,10
 268:	00078513          	mv	a0,a5
 26c:	00000097          	auipc	ra,0x0
 270:	000080e7          	jalr	ra # 26c <.L13+0x4c>
 274:	00050793          	mv	a5,a0
 278:	fef42623          	sw	a5,-20(s0)
 27c:	fec42783          	lw	a5,-20(s0)
 280:	fa0790e3          	bnez	a5,220 <.L13>
 284:	fd842783          	lw	a5,-40(s0)
 288:	0207d263          	bgez	a5,2ac <.L14>
 28c:	fe842783          	lw	a5,-24(s0)
 290:	00178713          	addi	a4,a5,1
 294:	fee42423          	sw	a4,-24(s0)
 298:	00078713          	mv	a4,a5
 29c:	fdc42783          	lw	a5,-36(s0)
 2a0:	00e787b3          	add	a5,a5,a4
 2a4:	02d00713          	li	a4,45
 2a8:	00e78023          	sb	a4,0(a5)

000002ac <.L14>:
 2ac:	fe842783          	lw	a5,-24(s0)
 2b0:	fdc42703          	lw	a4,-36(s0)
 2b4:	00f707b3          	add	a5,a4,a5
 2b8:	00078023          	sb	zero,0(a5)
 2bc:	fe842783          	lw	a5,-24(s0)
 2c0:	fff78793          	addi	a5,a5,-1
 2c4:	fef42423          	sw	a5,-24(s0)
 2c8:	0640006f          	j	32c <.L15>

000002cc <.L16>:
 2cc:	fe442783          	lw	a5,-28(s0)
 2d0:	fdc42703          	lw	a4,-36(s0)
 2d4:	00f707b3          	add	a5,a4,a5
 2d8:	0007c783          	lbu	a5,0(a5)
 2dc:	fef401a3          	sb	a5,-29(s0)
 2e0:	fe842783          	lw	a5,-24(s0)
 2e4:	fdc42703          	lw	a4,-36(s0)
 2e8:	00f70733          	add	a4,a4,a5
 2ec:	fe442783          	lw	a5,-28(s0)
 2f0:	00178693          	addi	a3,a5,1
 2f4:	fed42223          	sw	a3,-28(s0)
 2f8:	00078693          	mv	a3,a5
 2fc:	fdc42783          	lw	a5,-36(s0)
 300:	00d787b3          	add	a5,a5,a3
 304:	00074703          	lbu	a4,0(a4)
 308:	00e78023          	sb	a4,0(a5)
 30c:	fe842783          	lw	a5,-24(s0)
 310:	fff78713          	addi	a4,a5,-1
 314:	fee42423          	sw	a4,-24(s0)
 318:	00078713          	mv	a4,a5
 31c:	fdc42783          	lw	a5,-36(s0)
 320:	00e787b3          	add	a5,a5,a4
 324:	fe344703          	lbu	a4,-29(s0)
 328:	00e78023          	sb	a4,0(a5)

0000032c <.L15>:
 32c:	fe442703          	lw	a4,-28(s0)
 330:	fe842783          	lw	a5,-24(s0)
 334:	f8f74ce3          	blt	a4,a5,2cc <.L16>
 338:	fdc42783          	lw	a5,-36(s0)
 33c:	00078513          	mv	a0,a5
 340:	02c12083          	lw	ra,44(sp)
 344:	02812403          	lw	s0,40(sp)
 348:	03010113          	addi	sp,sp,48
 34c:	00008067          	ret

00000350 <main>:
 350:	fd010113          	addi	sp,sp,-48
 354:	02112623          	sw	ra,44(sp)
 358:	02812423          	sw	s0,40(sp)
 35c:	03010413          	addi	s0,sp,48
 360:	fd040793          	addi	a5,s0,-48
 364:	00000593          	li	a1,0
 368:	00078513          	mv	a0,a5
 36c:	00000097          	auipc	ra,0x0
 370:	000080e7          	jalr	ra # 36c <main+0x1c>
 374:	00050793          	mv	a5,a0
 378:	00078513          	mv	a0,a5
 37c:	00000097          	auipc	ra,0x0
 380:	000080e7          	jalr	ra # 37c <main+0x2c>
 384:	fd040793          	addi	a5,s0,-48
 388:	4d200593          	li	a1,1234
 38c:	00078513          	mv	a0,a5
 390:	00000097          	auipc	ra,0x0
 394:	000080e7          	jalr	ra # 390 <main+0x40>
 398:	00050793          	mv	a5,a0
 39c:	00078513          	mv	a0,a5
 3a0:	00000097          	auipc	ra,0x0
 3a4:	000080e7          	jalr	ra # 3a0 <main+0x50>
 3a8:	fd040793          	addi	a5,s0,-48
 3ac:	b2e00593          	li	a1,-1234
 3b0:	00078513          	mv	a0,a5
 3b4:	00000097          	auipc	ra,0x0
 3b8:	000080e7          	jalr	ra # 3b4 <main+0x64>
 3bc:	00050793          	mv	a5,a0
 3c0:	00078513          	mv	a0,a5
 3c4:	00000097          	auipc	ra,0x0
 3c8:	000080e7          	jalr	ra # 3c4 <main+0x74>
 3cc:	fd040713          	addi	a4,s0,-48
 3d0:	800007b7          	lui	a5,0x80000
 3d4:	fff7c593          	not	a1,a5
 3d8:	00070513          	mv	a0,a4
 3dc:	00000097          	auipc	ra,0x0
 3e0:	000080e7          	jalr	ra # 3dc <main+0x8c>
 3e4:	00050793          	mv	a5,a0
 3e8:	00078513          	mv	a0,a5
 3ec:	00000097          	auipc	ra,0x0
 3f0:	000080e7          	jalr	ra # 3ec <main+0x9c>
 3f4:	fd040793          	addi	a5,s0,-48
 3f8:	800005b7          	lui	a1,0x80000
 3fc:	00078513          	mv	a0,a5
 400:	00000097          	auipc	ra,0x0
 404:	000080e7          	jalr	ra # 400 <main+0xb0>
 408:	00050793          	mv	a5,a0
 40c:	00078513          	mv	a0,a5
 410:	00000097          	auipc	ra,0x0
 414:	000080e7          	jalr	ra # 410 <main+0xc0>
 418:	fd040793          	addi	a5,s0,-48
 41c:	0ab00593          	li	a1,171
 420:	00078513          	mv	a0,a5
 424:	00000097          	auipc	ra,0x0
 428:	000080e7          	jalr	ra # 424 <main+0xd4>
 42c:	00050793          	mv	a5,a0
 430:	00078513          	mv	a0,a5
 434:	00000097          	auipc	ra,0x0
 438:	000080e7          	jalr	ra # 434 <main+0xe4>
 43c:	fd040713          	addi	a4,s0,-48
 440:	0000b7b7          	lui	a5,0xb
 444:	bcd78593          	addi	a1,a5,-1075 # abcd <main+0xa87d>
 448:	00070513          	mv	a0,a4
 44c:	00000097          	auipc	ra,0x0
 450:	000080e7          	jalr	ra # 44c <main+0xfc>
 454:	00050793          	mv	a5,a0
 458:	00078513          	mv	a0,a5
 45c:	00000097          	auipc	ra,0x0
 460:	000080e7          	jalr	ra # 45c <main+0x10c>
 464:	fd040713          	addi	a4,s0,-48
 468:	deadc7b7          	lui	a5,0xdeadc
 46c:	eef78593          	addi	a1,a5,-273 # deadbeef <main+0xdeadbb9f>
 470:	00070513          	mv	a0,a4
 474:	00000097          	auipc	ra,0x0
 478:	000080e7          	jalr	ra # 474 <main+0x124>
 47c:	00050793          	mv	a5,a0
 480:	00078513          	mv	a0,a5
 484:	00000097          	auipc	ra,0x0
 488:	000080e7          	jalr	ra # 484 <main+0x134>
 48c:	00000793          	li	a5,0
 490:	00078513          	mv	a0,a5
 494:	02c12083          	lw	ra,44(sp)
 498:	02812403          	lw	s0,40(sp)
 49c:	03010113          	addi	sp,sp,48
 4a0:	00008067          	ret

*/

//-----------------------------------------------------------------------------

func main() {



  os.Exit(0)
}

//-----------------------------------------------------------------------------

