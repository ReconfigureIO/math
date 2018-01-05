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

func fixed26(us <-chan uint32) fixed.Int26_6 {
	return fixed.I26F(0, int32(<-us))
}

// restricted ln from [0, 1) using a 32 lookup table
func log(x fixed.Int26_6) fixed.Int26_6 {
	return [32]fixed.Int26_6{0, 1, 3, 5, 7, 9, 10, 12, 14, 15, 17, 18, 20, 21, 23, 24, 25, 27, 28, 29, 31, 32, 33, 34, 35, 36, 38, 39, 40, 41, 42, 43}[(x>>1)&0x1f]
}

func normals(uint32s <-chan uint32) fixed.Int26_6 {
	xs := [c]fixed.Int26_6{0, 13, 18, 21, 24, 26, 28, 29, 31, 32, 33, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 44, 45, 46, 47, 47, 48, 49, 50, 50, 51, 52, 52, 53, 54, 54, 55, 55, 56, 57, 57, 58, 58, 59, 59, 60, 61, 61, 62, 62, 63, 63, 64, 64, 65, 65, 66, 66, 67, 67, 68, 68, 69, 69, 70, 70, 71, 71, 72, 72, 72, 73, 73, 74, 74, 75, 75, 76, 76, 77, 77, 77, 78, 78, 79, 79, 80, 80, 81, 81, 81, 82, 82, 83, 83, 84, 84, 84, 85, 85, 86, 86, 87, 87, 87, 88, 88, 89, 89, 90, 90, 90, 91, 91, 92, 92, 93, 93, 93, 94, 94, 95, 95, 96, 96, 96, 97, 97, 98, 98, 99, 99, 100, 100, 100, 101, 101, 102, 102, 103, 103, 103, 104, 104, 105, 105, 106, 106, 107, 107, 108, 108, 108, 109, 109, 110, 110, 111, 111, 112, 112, 113, 113, 114, 114, 115, 115, 116, 116, 117, 117, 118, 118, 119, 119, 120, 120, 121, 121, 122, 122, 123, 123, 124, 125, 125, 126, 126, 127, 127, 128, 129, 129, 130, 130, 131, 132, 132, 133, 133, 134, 135, 135, 136, 137, 137, 138, 139, 140, 140, 141, 142, 143, 143, 144, 145, 146, 147, 147, 148, 149, 150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 163, 164, 165, 167, 168, 170, 171, 173, 174, 176, 178, 180, 182, 185, 187, 190, 193, 197, 201, 206, 212, 220, 233}

	fs := [c]fixed.Int26_6{64, 62, 61, 60, 59, 58, 58, 57, 56, 56, 55, 55, 54, 53, 53, 52, 52, 51, 51, 50, 50, 50, 49, 49, 48, 48, 47, 47, 47, 46, 46, 45, 45, 45, 44, 44, 44, 43, 43, 42, 42, 42, 41, 41, 41, 40, 40, 40, 39, 39, 39, 39, 38, 38, 38, 37, 37, 37, 36, 36, 36, 35, 35, 35, 35, 34, 34, 34, 33, 33, 33, 33, 32, 32, 32, 32, 31, 31, 31, 30, 30, 30, 30, 29, 29, 29, 29, 28, 28, 28, 28, 27, 27, 27, 27, 26, 26, 26, 26, 26, 25, 25, 25, 25, 24, 24, 24, 24, 23, 23, 23, 23, 23, 22, 22, 22, 22, 22, 21, 21, 21, 21, 20, 20, 20, 20, 20, 19, 19, 19, 19, 19, 18, 18, 18, 18, 18, 17, 17, 17, 17, 17, 16, 16, 16, 16, 16, 15, 15, 15, 15, 15, 15, 14, 14, 14, 14, 14, 13, 13, 13, 13, 13, 13, 12, 12, 12, 12, 12, 11, 11, 11, 11, 11, 11, 10, 10, 10, 10, 10, 10, 9, 9, 9, 9, 9, 9, 9, 8, 8, 8, 8, 8, 8, 7, 7, 7, 7, 7, 7, 7, 6, 6, 6, 6, 6, 6, 5, 5, 5, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4, 4, 3, 3, 3, 3, 3, 3, 3, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	ms := [c]fixed.Int26_6{0, -14, -18, -20, -22, -24, -25, -26, -27, -28, -29, -30, -31, -31, -32, -32, -33, -33, -34, -34, -34, -35, -35, -35, -36, -36, -36, -36, -36, -37, -37, -37, -37, -37, -37, -37, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -38, -37, -37, -37, -37, -37, -37, -37, -37, -37, -37, -36, -36, -36, -36, -36, -36, -36, -36, -36, -35, -35, -35, -35, -35, -35, -35, -35, -34, -34, -34, -34, -34, -34, -34, -33, -33, -33, -33, -33, -33, -32, -32, -32, -32, -32, -32, -31, -31, -31, -31, -31, -31, -30, -30, -30, -30, -30, -30, -29, -29, -29, -29, -29, -28, -28, -28, -28, -28, -27, -27, -27, -27, -27, -26, -26, -26, -26, -26, -25, -25, -25, -25, -25, -24, -24, -24, -24, -24, -23, -23, -23, -23, -22, -22, -22, -22, -22, -21, -21, -21, -21, -20, -20, -20, -20, -20, -19, -19, -19, -19, -18, -18, -18, -18, -18, -17, -17, -17, -17, -16, -16, -16, -16, -15, -15, -15, -15, -14, -14, -14, -14, -13, -13, -13, -13, -12, -12, -12, -12, -11, -11, -11, -11, -10, -10, -10, -10, -9, -9, -9, -9, -8, -8, -8, -8, -7, -7, -7, -7, -6, -6, -6, -5, -5, -5, -5, -4, -4, -4, -4, -3, -3, -3, -3, -2, -2, -2, -1, -1, -1, -1, 0, 0, 0}

	keepGoing := true
	var out fixed.Int26_6
	for keepGoing {
		u := int32(<-uint32s)
		// the index we'll use
		i := u & mask
		x := xs[i]
		// use u as a fixed point from [0..1)
		z := fixed.I26F(0, u).Mul(x)
		if i != c-1 && z < xs[i+1] {
			keepGoing = false
			// in bulk, this path should happen very frequently
			if u < 0 {
				out = -1 * z
			} else {
				out = z
			}
		} else if i == 0 {
			// Tail
			var x2 fixed.Int26_6
			for keepGoing {
				x2 := -log(fixed26(uint32s)).Mul(rInv)
				y := -log(fixed26(uint32s)) << 1
				if y >= x2*x2 {
					keepGoing = false
				}
			}
			if u > 0 {
				out = r + x2
			} else {
				out = -r - x2
			}
		} else {
			// wedge
			f := fs[i-1] - fs[i]
			if f < 0 {
				f = -1 * f
			}
			y := fixed26(uint32s).Mul(f)
			if y < ms[i-1]*(z-x) {
				keepGoing = false
				if u < 0 {
					out = -1 * y
				} else {
					out = y
				}
			}
		}

	}
	return out
}

// Normals writes a stream of Int26_6, normally distributed
func (r Rand) Normals(output chan<- fixed.Int26_6) {
	uint32s := make(chan uint32, 1)
	go r.Uint32s(uint32s)

	for {
		output <- normals(uint32s)
	}

}
