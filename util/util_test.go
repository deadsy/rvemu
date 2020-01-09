//-----------------------------------------------------------------------------
/*

Utility Function Testing

*/
//-----------------------------------------------------------------------------

package util

import (
	"fmt"
	"testing"
)

//-----------------------------------------------------------------------------

type u64Test struct {
	a, b   uint64 // operands
	result uint64 // expected result
}

type s64Test struct {
	a, b   int64 // operands
	result int64 // expected result
}

//-----------------------------------------------------------------------------

func Test_Mulhu(t *testing.T) {

	test := []u64Test{
		{0, 0, 0},
		{1, 1, 0},
		{3, 7, 0},
		{0xffff8000, 0, 0},
		{0x80000000, 0, 0},
		{0, 0xffffffffffff8000, 0},
		{0xffffffff80000000, 0, 0},
		{0xffffffff80000000, 0xffffffffffff8000, 0xffffffff7fff8000},
		{0xaaaaaaaaaaaaaaab, 0x2fe7d, 0x1fefe},
		{0x2fe7d, 0xaaaaaaaaaaaaaaab, 0x1fefe},
		{1 << 32, 1 << 32, 1},
		{1 << 31, 1 << 32, 0},
		{1 << 33, 1 << 32, 1 << 1},
		{1 << 63, 1 << 63, 1 << 62},
		{0xffffffffffffffff, 0xffffffffffffffff, 0xfffffffffffffffe},
		{0xea7b9ab1dd6322de, 0x0aa795129ef02221, 0x09c25337dbfbee02},
		{0x04f9a24cfaf8f7cc, 0x3f3d9bf3638a7fd4, 0x013aa1747e43f826},
		{0xd461d4f02995033e, 0x048bc6fe3852ec09, 0x03c57d88f18a2024},
		{0xe2cbb8fd5202469b, 0x71688939af234570, 0x647888fd18a09635},
		{0x4cf95b8423c2b7ab, 0x738550450edb70c2, 0x22bc19cb97fda0af},
		{0x4f6634567cb1f955, 0x4033cb02768c2aef, 0x13e99d68db78e8ad},
		{0x9592446ada52c58a, 0x04439adc8a120cdb, 0x027dc8d46ae6e696},
		{0x47ffc47882fff842, 0x6d1f052a4174916b, 0x1eb0a013fb9429b5},
		{0xff74a2990abef66b, 0xf096c576dee897c5, 0xf013cbd21f1148df},
		{0x2042e1e79ad60ec6, 0x22007a90d398cabf, 0x0448f172e282417c},
		{0x77016ce93caa5fbb, 0xf6a1e0a8a4efd182, 0x72a69efd43e2d2c7},
		{0xdf1d5a174ae711ea, 0x8414e23ae68483a8, 0x731d55ea4b87d582},
		{0x7c6a7d0b3690606b, 0xddf06544f222da83, 0x6bdcc2ff6dea08a3},
		{0xa491d13354d88dc4, 0x5a2b024019b001d1, 0x39f6e97c880d788e},
		{0xecd885c9abd5db23, 0x518fbf0a15fb165f, 0x4b7580066ba267d7},
		{0x6d7ad81bfb04a6e6, 0xe718ff4b3d11e1fc, 0x62d489b30c7fd591},
		{0xbba855cf23368bc9, 0xa683560038aeacd6, 0x7a0f6dd2f034d743},
		{0x217d891650dc46a1, 0x4bf3515dcafb1f22, 0x09ef9bf7a7680ecc},
		{0xfcf30c4f6697e397, 0x00a434c14a1f5fa2, 0x00a23fd43dd7637c},
		{0x534c59129532f8e5, 0xfb98c723fae939e9, 0x51dd917629a3a222},
	}

	for _, v := range test {
		x := Mulhu2(v.a, v.b)
		if x != v.result {
			fmt.Printf("%x x %x = %x (expected %x)\n", v.a, v.b, x, v.result)
			t.Error("FAIL")
		}
		x = Mulhu2(v.b, v.a)
		if x != v.result {
			fmt.Printf("%x x %x = %x (expected %x)\n", v.b, v.a, x, v.result)
			t.Error("FAIL")
		}
	}
}

func Test_Mulhs(t *testing.T) {

	test := []s64Test{
		{0, 0, 0},
		{1, 1, 0},
		{-1, -1, 0},
		{1, -1, -1},
		{10, -10, -1},
		{1 << 32, 1 << 32, 1},
		{1 << 31, 1 << 32, 0},
		{1 << 33, 1 << 32, 1 << 1},
		{1433146433078531771, -9068335344185109807, -704528257157382784},
		{4829535485415149908, -8922292195709108081, -2335942136901696229},
		{-230731395066329254, 3042306729403585405, -38053093439692895},
		{-6625085407879829436, -7108697370776146514, 2553064488344258425},
		{-3068985225564752887, 6153461497047533421, -1023751527373078973},
		{-7289289615376117957, 1496886750032238122, -591499562134376275},
		{-5610198934957303438, 8041075967822838614, -2445528362638504557},
		{3861963970238710803, 5328844735747868618, 1115633538917309444},
		{-7157475940990982760, 5071531100968786240, -1967792348293363123},
		{5215402259787899221, 7238409379830374187, 2046497554581493761},
		{-7766700907006896043, 3555347232371209510, -1496920999393961060},
		{-1638696414469481542, -8307099383434620829, 737952124227551149},
		{8307650744875703509, -1534433469992331852, -691045384974571160},
		{-5119214265624239148, 5053967338379839003, -1402542453738842080},
		{2163596123780717888, -683427565732920843, -80158386010897134},
		{-8888281125798942958, 3566902160659143292, -1718657180122230900},
		{1034952864523846726, 6720781308292035886, 377068811659206401},
		{-7000632227592065292, 4033654152470200075, -1530792054245965846},
		{8507963851107869321, 8435934281062632489, 3890801738605646503},
		{-651759051721940080, -2000325471223497852, 70675357507559111},
	}

	for _, v := range test {
		x := Mulhs(v.a, v.b)
		if x != v.result {
			fmt.Printf("%x x %x = %x (expected %x)\n", v.a, v.b, x, v.result)
			t.Error("FAIL")
		}
		x = Mulhs(v.b, v.a)
		if x != v.result {
			fmt.Printf("%x x %x = %x (expected %x)\n", v.b, v.a, x, v.result)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
