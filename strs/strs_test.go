package strs

import (
	"testing"
)

func specEq(t *testing.T, a, b Spec) {
	if a != b {
		t.Fail()
	}
}

func TestSpec(t *testing.T) {
	specEq(t, NoSymbol, ParseSpec("lun"))
	specEq(t, Num, ParseSpec("n"))
	specEq(t, Num, ParseSpec("num"))
	specEq(t, Lower, ParseSpec("lower"))
	specEq(t, Upper, ParseSpec("upper"))
	specEq(t, Symbol, ParseSpec("symbol"))
	specEq(t, All, ParseSpec("luns"))
	specEq(t, NoSymbol, ParseSpec("abc"))
	specEq(t, Upper|Lower, ParseSpec("upper", "lower", "l", "x"))
}

func TestRandString(t *testing.T) {
	_ = All.String()
	w := map[Spec]int{
		Num:    10,
		Symbol: 18,
		Upper:  26,
		Lower:  26,
	}

	RandString(0, 0, false)
	RandString(1, 0, false)
	RandString(2, NoSymbol, false)
	for spec, l := range w {
		RandString(l/2, spec, true)
		RandString(l/2, spec, false)
		RandString(l-1, spec, true)
		RandString(l-1, spec, false)
		RandString(l+1, spec, true)
		RandString(l+1, spec, false)
		RandString(l, spec, true)
		RandString(l, spec, false)
	}
}

func BenchmarkRandStringNoSymbol(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(12, NoSymbol, false)
	}
}

func BenchmarkRandStringAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(12, All, false)
	}
}

func BenchmarkRandStringUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandString(12, Upper, false)
	}
}

func TestRandHs(t *testing.T) {
	RandHs(12, false)
	RandHs(10901, true)
	RandHs(20903, true)
}

func BenchmarkRandHs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandHs(12, true)
	}
}

func TestRandSlice(t *testing.T) {
	RandSlice([]interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	RandSlice2([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	RandSlice2("si")
	RandSlice2(nil)
}

func BenchmarkRandSlice(b *testing.B) {
	var s = []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < b.N; i++ {
		RandSlice(s)
	}
}

func BenchmarkRandSlice2(b *testing.B) {
	var s = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i < b.N; i++ {
		RandSlice2(s)
	}
}

func TestIncr(t *testing.T) {
	for i := 0; i < 100000000; i++ {
		Incr()
	}
	t.Log(Incrs())
	t.Log(Incrs())
	t.Log(Incrs())
	t.Log(Incrs())
	t.Log(Incrs())
}

func BenchmarkRandIncr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Incr()
	}
	b.Log(Incrs())
}
