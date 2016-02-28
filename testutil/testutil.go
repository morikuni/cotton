package testutil

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"testing"
)

func NOPHandler(w http.ResponseWriter, r *http.Request) {}

type T struct {
	*testing.T
}

func (t *T) println(s string) {
	fmt.Fprintln(os.Stderr, "\t\t"+s)
}

func (t *T) printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "\t\t"+format, args...)
}

func (t *T) fail(name string) {
	var file string
	var line int
	_, file, line, ok := runtime.Caller(2)
	if ok {
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	fmt.Fprintf(os.Stderr, "\t%s:%d:%s\n", file, line, name)
	t.Fail()
}

func (t *T) MustEqual(l interface{}, r interface{}) {
	if l != r {
		t.fail("MustEqual")
		t.printf(" Left: %#v\n", l)
		t.printf("Right: %#v\n", r)
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
		cm.t.fail("MustCalled")
	}
}

func (cm *callMe) MustCalledTimes(n uint) {
	if cm.count != n {
		cm.t.fail("MustCalledTimes")
		cm.t.printf("Expect: %d times\n", n)
		cm.t.printf("Actual: %d times\n", cm.count)
	}
}
