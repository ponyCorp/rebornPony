package cmdhandler

import (
	"fmt"
	"testing"
)

func TestAddRules(t *testing.T) {

	r := newRules()
	r.add("test1")
	r.add("test2")
	r.add("test3")
	fmt.Println(r.get())
}
