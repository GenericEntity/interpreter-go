package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	t1 := &Boolean{Value: true}
	t2 := &Boolean{Value: true}
	f1 := &Boolean{Value: false}
	f2 := &Boolean{Value: false}

	if t1.HashKey() != t2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}
	if f1.HashKey() != f2.HashKey() {
		t.Errorf("booleans with same content have different hash keys")
	}
	if t1.HashKey() == f1.HashKey() {
		t.Errorf("booleans with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	a1 := &Integer{Value: 1}
	a2 := &Integer{Value: 1}
	b1 := &Integer{Value: 2}
	b2 := &Integer{Value: 2}

	if a1.HashKey() != a2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}
	if b1.HashKey() != b2.HashKey() {
		t.Errorf("integers with same content have different hash keys")
	}
	if a1.HashKey() == b1.HashKey() {
		t.Errorf("integers with different content have same hash keys")
	}
}
