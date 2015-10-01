package cash

import (
	"testing"
	"time"
)

const KEY = "foo"
const VAL = "bar"

var res interface{}

func BenchmarkGet(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	c.Set(KEY, VAL, Default)

	var str string

	for n := 0; n < b.N; n++ {
		if x, y := c.Get(KEY); y {
			str = x.(string)
		}
	}

	res = str
}

func BenchmarkSet(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	for n := 0; n < b.N; n++ {
		c.Set(KEY, VAL, Default)
	}
}

func BenchmarkHas(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	c.Set(KEY, VAL, Default)

	var has bool

	for n := 0; n < b.N; n++ {
		res = c.Has(KEY)
	}

	res = has
}

func BenchmarkDel(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	c.Set(KEY, VAL, Default)

	for n := 0; n < b.N; n++ {
		c.Del(KEY)
	}
}

func BenchmarkClean(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	c.Set(KEY, VAL, Default)

	for n := 0; n < b.N; n++ {
		c.Clean()
	}
}

func BenchmarkFlush(b *testing.B) {
	c := New(Conf{10 * time.Minute, -1})

	c.Set(KEY, VAL, Default)

	for n := 0; n < b.N; n++ {
		c.Flush()
	}
}
