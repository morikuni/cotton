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

func (t *T) CallMe() *CallMe {
	return &CallMe{
		t:    t,
		flag: false,
	}
}

type CallMe struct {
	t    *T
	flag bool
}

func (cm *CallMe) Call() {
	cm.flag = true
}

func (cm *CallMe) MustCalled() {
	if cm.flag == false {
		cm.t.Error("Error: MustCalled")
	}
}
