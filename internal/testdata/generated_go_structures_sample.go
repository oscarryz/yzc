package main

import "fmt"

/*
// filename is e.yz
s : "The sound of the ocean calm my soul"

g : {
	t String
}

f : {
	s : "The sound of the ocean calm my soul"
	g : {
		t String
		println(t)
		t.reverse()
	}
	h : {
		s.reverse()
	}
}

print("Hello `s`, `g`, `f`, `f.s`, `f.g`, `f.g.t`")
s = f.g("This is fun")
*/
type (
	_e struct {
		s  string
		g  *_e_g
		_f *_e_f
	}
	_e_f struct {
		s string
		g *_e_f_g
		h *_e_f_h
	}
	_e_f_g struct {
		t string
	}
	_e_f_h struct {
		__f *_e_f
		results []any
	}
	_e_g struct {
		t string
	}
)

func (gfg *_e_f_g) apply() {
	fmt.Printf("%s\n", gfg.t)
	gfg.t = std_string_reverse(gfg.t)
}
func (efh *_e_f_h) apply() {
	std_string_reverse(efh.__f.s)
}

func std_string_reverse(s string) string {
	return s
}
func (f *_e_f) apply() {
}
func (e *_e) apply() {
	e._f = &_e_f{
		s: "The sound of the ocean calm my soul",
		g: &_e_f_g{},
		h: &_e_f_h{},
	}
	e._f.h.__f = e._f


	fmt.Printf("Hello %s, %v, %v, %s, %v, %v", e.s, e.g, e._f, e._f.s, e._f.g, e._f.g.t)
	e._f.g.t = "This is fun"
	e._f.g.apply()
	e.s = e._f.g.t

}
func main() {
	e := _e{}
	e.apply()
	f:= &_factorial{}
	f.n = 4
	f.apply()
	fmt.Println(f.result)

}

/*
factorial #(n Int, Int)
factoria = { n Int
	n == 0 ? { 1 }
	{ n * factorial(n -1 )}
}
factorial(2)
*/
type (
	_factorial struct { 
		n int
		result int
	}
)
func (f *_factorial) apply() {
	if f.n == 0 { 
		f.result = 1
		return 
	} else {
		n := f.n
		f.n = f.n - 1
		f.apply()
		f.result = n * f.result
	}
}