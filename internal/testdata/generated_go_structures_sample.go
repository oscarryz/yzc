package main

import (
	"fmt"
	"sync"
)

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
	// conc helps to run functions concurrently
	// and waits until they complete
	// similar idea as https://github.com/sourcegraph/conc but simplified for this sample
	conc struct {
		wg sync.WaitGroup
	}
	runner interface {
		run()
	}
	// transformed structs
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
		__f     *_e_f
		results []any
	}
	_e_g struct {
		t string
	}
	_factorial struct {
		n      int
		result int
	}
)

func (c *conc) Run(r runner) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		r.run()
	}()
}
func (gfg *_e_f_g) run() {
	fmt.Printf("%s\n", gfg.t)
	gfg.t = std_string_reverse(gfg.t)
}
func (efh *_e_f_h) run() {
	std_string_reverse(efh.__f.s)
}

func std_string_reverse(s string) string {
	return s
}
func (f *_e_f) run() {
}
func (e *_e) run() {
	e._f = &_e_f{
		s: "The sound of the ocean calm my soul",
		g: &_e_f_g{},
		h: &_e_f_h{},
	}
	e._f.h.__f = e._f

	fmt.Printf("Hello %s, %v, %v, %s, %v, %v\n", e.s, e.g, e._f, e._f.s, e._f.g, e._f.g.t)
	e._f.g.t = "This is fun"
	e._f.g.run()
	e.s = e._f.g.t

}

/*
factorial #(n Int, Int)
factoria = { n Int

		n == 0 ? { 1 }
		{ n * factorial(n -1 )}
	}

factorial(2)
*/
func (f *_factorial) run() {
	if f.n == 0 {
		f.result = 1
		return
	} else {
		n := f.n
		f.n = f.n - 1
		f.run()
		f.result = n * f.result
	}
}
func main() {
	var c conc
	e := &_e{}
	c.Run(e)
	f := &_factorial{}
	f.n = 4
	c.Run(f)
	c.wg.Wait()
	fmt.Println(f.result)

}
