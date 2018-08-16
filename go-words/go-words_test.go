// go-words project words.go
package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	a := assert.New(t)
	a.Equal([]int{3, 4, 6, 5, 7}, unique([]int{3, 4, 6}, []int{3, 5, 6, 7}))
}

func TestProcessorNext(t *testing.T) {
	//	a := assert.New(t)
	cl := NewProcessor()
	//phrase := "<aabbccaabbccaabbccaaccaabbcc>"
	//phrase := "abc_ac_abc_ac_abc_"

	phrase := "abc_abc_abcd_abc_abcd_abc_"
	// [abc_] x 4 = 16, [abcd_] x 2 = 10

	// [abc] x 6 = 18, [_][d][c]
	// [_abc] x 5 = 20, [abc][d][_]

	// [abc_abcd_]

	//phrase := "abc_l_abc_i"

	for _, b := range []byte(phrase) {
		cl.Next(Char(b))
	}
	cl.Done()

	//	cl.Clean()
	//	cl.SortPositions()

	for _, tu := range cl.Vocabulary() {
		t.Log(tu)
	}
}
