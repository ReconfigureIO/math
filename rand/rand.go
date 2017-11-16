package rand

func RandUint32(seed uint32, output chan<- uint32) {
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
