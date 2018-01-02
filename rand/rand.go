package rand

// A Rand is a source of random numbers
type Rand struct {
	seed uint32
}

// Construct a new Rand given a seed
func New(seed uint32) Rand {
	return Rand{seed: seed}
}

// Write a stream of uint32s to the given channel
func (r Rand) Uint32s(output chan<- uint32) {
	seed := r.seed
	for {
		a := seed ^ (seed << 13)
		b := a ^ (a >> 17)
		c := b ^ (b << 5)
		go func() {
			output <- c
		}()
		seed = c
	}
}
