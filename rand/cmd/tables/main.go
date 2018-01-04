package main

import (
	"fmt"
	"math"

	"github.com/ReconfigureIO/fixed/host"
)

const (
	c = 256
	r = 3.442619855899
)

func f(x float64) float64 {
	return math.Exp(-x * x / 2)
}

func fInv(y float64) float64 {
	return math.Sqrt(-2 * math.Log(y))
}

func printVar(varName string, vals []float64) {
	fmt.Printf("%s = {", varName)
	for i, v := range vals {
		if i != c-1 {
			fmt.Printf("%d, ", int32(host.I26Float64(v)))
		} else {
			fmt.Printf("%d}\n", int32(host.I26Float64(v)))
		}
	}

}

func main() {
	a := float64(0)
	b := float64(10)
	r := float64(5)
	aa := a
	bb := b
	xs := make([]float64, c, c)
	for aa < r && r < bb {
		aa = a
		bb = b
		r = .5 * (a + b)

		v := r*f(r) + math.Exp(-0.5*r*r)/r
		xs[c-1] = r
		for i := c - 1; i != 0; i-- {
			xs[i-1] = fInv(v/xs[i] + f(xs[i]))
		}
		q := xs[1] * (1 - f(xs[1])) / v
		if q > 1 {
			b = r
		} else {
			a = r
		}
	}
	xs[0] = 0

	ys := make([]float64, c, c)
	for i, x := range xs {
		ys[i] = f(x)
	}

	bs := make([]float64, c, c)
	for i := 1; i < c; i++ {
		dy := ys[i] - ys[i-1]
		bs[i] = ((2 / dy) * (dy * fInv((ys[i]+ys[i-1])/2))) - xs[i-1]
	}

	ms := make([]float64, c, c)
	for i := 1; i < c-1; i++ {
		ms[i] = (ys[i] - ys[i+1]) / (xs[i] - bs[i+1])
	}

	printVar("xs", xs)
	printVar("ys", ys)
	printVar("ms", ms)
	fmt.Printf("r = %d\n", host.I26Float64(r))
	fmt.Printf("rInv = %d\n", host.I26Float64(1.0/r))

}
