package main

import "github.com/ReconfigureIO/math/rand"

func Top() {
	channel := make(chan uint32)
	println(RandUint32(432459, channel))
}
