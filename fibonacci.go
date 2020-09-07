package main

import (
	"math"
)

//Provide for fibonacci sequence operations
type Tfibonacci struct {
	Cache map[uint]uint64
}

func (g *Tfibonacci) Init() {
	g.Cache = make(map[uint]uint64)
}

func (g *Tfibonacci) Base(idx uint) uint64 {
	//O(n) time, O(n) space
	if idx < 2 {
		return uint64(idx)
	}
	return g.Base(idx-2) + g.Base(idx-1)
}

func (g *Tfibonacci) Faster(idx uint) uint64 {
	//O(n) time, O(1) space
	var a uint64 = 0
	var b uint64 = 1
	var i uint
	for i = 0; i < idx; i++ {
		a, b = b, a+b
	}
	return a
}

func (g *Tfibonacci) Memoized(idx uint) uint64 {

	//Time - O(1) <-> O(n) time
	//Space - O(1)
	//O(1)+ time, O(1) space

	//Check the cache first
	if g.Cache == nil {
		g.Cache = make(map[uint]uint64)
	}

	fibo, bexists := g.Cache[idx]
	if bexists {
		return fibo
	}

	if idx < 2 {
		g.Cache[idx] = uint64(idx)
		return uint64(idx)
	}

	fibo = g.Memoized(idx-1) + g.Memoized(idx-2)
	g.Cache[idx] = fibo
	return fibo
}

//Approximation method instead of integer operations
func (g *Tfibonacci) Closed(idx uint) uint64 {
	//O(1) time, O(1) space
	//See: https://en.wikipedia.org/wiki/Fibonacci_number#Closed-form_expression
	fidx := float64(idx)
	phi := (1.0 + math.Sqrt(5.0)) / 2.0
	psi := -(1.0 / phi)
	return (uint64)((math.Pow(phi, fidx) - math.Pow(psi, fidx)) / (phi - psi))
}
