package stream

import (
	"fmt"
	"testing"
)

var testDoc1 = `
	<html>
		<body>
			<div>
				<p>Hello 1</p>
			</div>
			<div>
				<p>Hello 2</p>
			</div>
			<div>
				<p>Hello 3</p>
			</div>
		</body>
	</html>
`

func TestStream(t *testing.T) {
	s := NewEngine()

	err := s.LoadString(testDoc1)
	//err := s.LoadFile("../fox.html")
	if err != nil {
		t.Error(err)
	}

	s.Run(0, 4)

	for _, m := range s.Storage {
		if len(m.Indexes) != len(m.Values) {
			t.Error("indexes not equal to values")
		}
		size := len(m.Values)
		for i, row := range m.Values {
			if len(row)+i+1 != size {
				t.Errorf("values row has incorrect size %v vs %v", len(row)+i+1, size)
			}
		}
	}
}

func TestStreamGroups(t *testing.T) {
	s := NewEngine()

	err := s.LoadString(testDoc1)
	if err != nil {
		t.Error(err)
	}

	matrix := s.Run(0, 4)
	for i, row := range matrix.Matrix {
		fmt.Println(i, row.String())
	}
}
