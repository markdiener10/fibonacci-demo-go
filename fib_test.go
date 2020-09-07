package main

import (
	"testing"
)

//                    0  1  2  3  4  5  6  7  8    9   10  11  12   13   14   15   16   17
var gcases = []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597}

func TestFibBase(t *testing.T) {
	gfib := Tfibonacci{}
	for i, v := range gcases {
		if gfib.Base(uint(i)) == v {
			continue
		}
		t.Error("TestFibBase", i, v)
	}
}

func TestFibFaster(t *testing.T) {
	gfib := Tfibonacci{}
	for i, v := range gcases {
		if gfib.Faster(uint(i)) == v {
			continue
		}
		t.Error("TestFibFaster", i, v, gfib.Faster(uint(i)))
	}
}

func TestFibMemoized(t *testing.T) {

	gfib := Tfibonacci{}
	for i, v := range gcases {
		if gfib.Memoized(uint(i)) == v {
			continue
		}
		t.Error("TestFibMemoized", i, v, gfib.Memoized(uint(i)))
	}

}

//Approximated Fibonacci means we have some slop in integer values,
//It returns values within +/- 1 of actual value (maybe float rounding inside?)
func TestFibClosed(t *testing.T) {

	gfib := Tfibonacci{}
	var fibo uint64
	for i, v := range gcases {
		//We have some rounding error so our accuracy is +/- 1
		fibo = gfib.Closed(uint(i))
		if fibo == 0 {
			continue
		}
		if fibo >= v-1 {
			if fibo <= v+1 {
				continue
			}
		}
		t.Error("TestFibClosed", i, v, fibo)
	}
}

func TestCompareFastAlgorithms(t *testing.T) {
	gfib := Tfibonacci{}
	gfib.Init()
	var i uint
	for i = 0; i < 600; i++ {
		if gfib.Faster(i) != gfib.Memoized(i) {
			t.Error(i)
		}
	}
}

//40 takes too long on standard hardware for base case
func BenchmarkBase(b *testing.B) {
	gfib := Tfibonacci{}
	gfib.Base(30)
}

//Optimized versions below
const FIBCOUNT = 600

func BenchmarkFaster(b *testing.B) {
	gfib := Tfibonacci{}
	gfib.Faster(FIBCOUNT)
}

func BenchmarkMemoized(b *testing.B) {
	gfib := Tfibonacci{}
	gfib.Memoized(FIBCOUNT)
}

func BenchmarkClosed(b *testing.B) {
	gfib := Tfibonacci{}
	gfib.Closed(FIBCOUNT)
}
