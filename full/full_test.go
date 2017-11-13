// classify project full_test.go
package full

import (
	"os"
	"testing"

	"github.com/olesho/classify"
	"golang.org/x/net/html"
)

func TestFullClassify(t *testing.T) {
	f, err := os.Open("../examples/BBC - Homepage.html")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	n, err := html.Parse(f)
	if err != nil {
		t.Error(err)
	}
	a := classify.NewArena(*n)
	c := NewFullClassificator(a)
	c.Run()
	t.Log(c.bags)
}
