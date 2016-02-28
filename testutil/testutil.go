package testutil

import (
	"net/http"
	"testing"
)

func NOPHandler(w http.ResponseWriter, r *http.Request) {}

type T struct {
	*testing.T
}

func (t *T) MustEqual(a interface{}, b interface{}) {
	if a != b {
		t.Error("Error: MustEqual")
		t.Errorf(" Left: %#v", a)
		t.Errorf("Right: %#v", b)
	}
}

func (t *T) CallMe() CallMe {
	return newCallMe(t)
}

type CallMe interface {
	Call()
	MustCalled()
	MustCalledTimes(n uint)
}

type callMe struct {
	t     *T
	count uint
}

func newCallMe(t *T) CallMe {
	return &callMe{
		t:     t,
		count: 0,
	}
}

func (cm *callMe) Call() {
	cm.count++
}

func (cm *callMe) MustCalled() {
	if cm.count == 0 {
		cm.t.Error("Error: MustCalled")
	}
}

func (cm *callMe) MustCalledTimes(n uint) {
	if cm.count < n {
		cm.t.Error("Error: MustCalledTimes")
		cm.t.Errorf("Expect: %d times", n)
		cm.t.Errorf("Actual: %d times", cm.count)
	}
}
