package rand

import (
	"github.com/ReconfigureIO/fixed"
)

const (
	c    = 256
	mask = c - 1
	r    = 220
	rInv = 18
)

// Modified Ziggurat Algorithm based upon the paper
// "Hardware-Optimized Ziggurat Algorithm for High-Speed Gaussian Random Number Generators"

// restricted ln from [0, 1) using a 32 lookup table
func log(x fixed.Int26_6) fixed.Int26_6 {
	return [32]fixed.Int26_6{0, 1, 3, 5, 7, 9, 10, 12, 14, 15, 17, 18, 20, 21, 23, 24, 25, 27, 28, 29, 31, 32, 33, 34, 35, 36, 38, 39, 40, 41, 42, 43}[(x>>1)&0x1f]
}

type param struct {
	X     uint8
	F     uint8
	M     int8
	XNext uint8
}

// Normals writes a stream of Int26_6, normally distributed
func (rand Rand) Normals(output chan<- fixed.Int26_6) {
	// an encoding of c params as uint32s. This results is better performance than just an array of param, but we need to unpack it
	// see cmd/tables/main.go for how to generate it
	params := [c]uint32{13, 219611154, 303165973, 353365528, 403565594, 437053980, 470542365, 487253791, 520807968, 537519393, 554296355, 587850532, 604561957, 621338918, 638116135, 654893096, 671604777, 688381738, 705158955, 721935916, 738713132, 738713133, 755490094, 772201775, 788978991, 788978736, 805755953, 822533170, 839310386, 839310387, 856087348, 872864564, 872864565, 889641782, 906353462, 906353463, 923130679, 923130424, 939907641, 956684857, 956684858, 973462074, 973462075, 990239291, 990239292, 1007016509, 1023793725, 1023793726, 1040570942, 1040570943, 1057348159, 1057348160, 1074125376, 1074059841, 1090837057, 1090837058, 1107614274, 1107614275, 1124391491, 1124391492, 1141168708, 1141168709, 1157945925, 1157945926, 1174723142, 1174723143, 1191500359, 1191500360, 1208277576, 1208277576, 1208277577, 1225054793, 1225054794, 1241832266, 1241832267, 1258609483, 1258609484, 1275386700, 1275386701, 1292163917, 1292163917, 1292163918, 1308941134, 1308941391, 1325718607, 1325718608, 1342495824, 1342495825, 1359207505, 1359207505, 1359207506, 1375984722, 1375984979, 1392762195, 1392762196, 1409539412, 1409539412, 1409539413, 1426316629, 1426316630, 1443094102, 1443094103, 1459871319, 1459871319, 1459871320, 1476648536, 1476648537, 1493426009, 1493426010, 1510203226, 1510203226, 1510203227, 1526980443, 1526980700, 1543757916, 1543757917, 1560535133, 1560535133, 1560535134, 1577312606, 1577312607, 1594089823, 1594089824, 1610867040, 1610867040, 1610867297, 1627644513, 1627644514, 1644421730, 1644421731, 1661198947, 1661199204, 1677976420, 1677976420, 1677976421, 1694753637, 1694753894, 1711531110, 1711531111, 1728308327, 1728308327, 1728308584, 1745085800, 1745085801, 1761863017, 1761863018, 1778640490, 1778640491, 1795417707, 1795417708, 1812129388, 1812129644, 1812129645, 1828906861, 1828906862, 1845684078, 1845684335, 1862461551, 1862461552, 1879238768, 1879238769, 1896016241, 1896016242, 1912793458, 1912793459, 1929570931, 1929570932, 1946348148, 1946348149, 1963125365, 1963125622, 1979902838, 1979902839, 1996680055, 1996680312, 2013457528, 2013457529, 2030234745, 2030234746, 2047012218, 2047012219, 2063789435, 2063789436, 2080566909, 2097344125, 2097344126, 2114121342, 2114121343, 2130898815, 2130898816, 2147676033, 2164453249, 2164453506, 2181230722, 2181230723, 2198007940, 2214785412, 2214785413, 2231562629, 2231562630, 2248340103, 2265117319, 2265117320, 2281894537, 2298672009, 2298672010, 2315449227, 2332226444, 2349003916, 2349003917, 2365781134, 2382558351, 2399335823, 2399335824, 2416113041, 2432890258, 2449667731, 2466444947, 2466444948, 2483222165, 2499999638, 2516776855, 2533554072, 2550331289, 2567108762, 2583885979, 2600663196, 2617440413, 2634217886, 2650995103, 2667772320, 2684549537, 2701261475, 2734815908, 2751593125, 2768370599, 2801925032, 2818702250, 2852256683, 2869034157, 2902588590, 2919365808, 2952920242, 2986474932, 3020029366, 3053583801, 3103915451, 3137470142, 3187801793, 3238133445, 3305242569, 3372351438, 3456237524, 3556900828, 3691053289, 3909157119}

	uint32s := make(chan uint32, 1)
	go rand.Uint32s(uint32s)

	for {
		keepGoing := true
		var out fixed.Int26_6
		for keepGoing {
			u := int32(<-uint32s)
			// the index we'll use
			i := u & mask

			// choose our param param and unpack them into a struct
			// Unpacking should be nearly free
			tmp := params[i]
			p := param{uint8(tmp >> 24), uint8(tmp >> 16), int8(tmp >> 8), uint8(tmp)}

			x := fixed.Int26_6(p.X)

			// use u as a fixed point from [0..1)
			t := fixed.I26F(0, int32(uint32(u)>>6))
			z := t.Mul(x)

			if z < fixed.Int26_6(p.XNext) {
				// Bulk, this path should happen very frequently
				keepGoing = false
				if u < 0 {
					out = -z
				} else {
					out = z
				}

			} else if i == 0 {
				// Tail
				var x2 fixed.Int26_6
				for keepGoing {
					t := -log(fixed.I26F(0, int32(<-uint32s)))
					x2 = t.Mul(rInv)
					y := -log(fixed.I26F(0, int32(<-uint32s))) << 1
					if y >= x2.Mul(x2) {
						keepGoing = false
					}
				}
				if u < 0 {
					out = -r - x2
				} else {
					out = r + x2
				}
			} else {
				// wedge

				// This is actually a 22.10
				f := fixed.Int26_6(p.F)
				t := fixed.I26F(0, int32(<-uint32s))

				// Resulting fixed point mult will be a (22 + 26).(10 + 6), so we shift 10 and then cast to get us back to 26.6
				y := fixed.Int26_6((uint64(f) * uint64(t)) >> 10)
				m := fixed.Int26_6(p.M)

				if y < m.Mul(z-x) {
					keepGoing = false
					if u < 0 {
						out = -y
					} else {
						out = y
					}
				}
			}

		}
		output <- out
	}

}
