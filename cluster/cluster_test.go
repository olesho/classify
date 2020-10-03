package cluster

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/olesho/classify/arena"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestRenderLabels(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("../bbc.html")
	a.NoError(err)
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	a.NoError(err)

	arena := arena.NewArena()
	arena.Append(*n)
	series := ExtractOptimized(arena).Matrix[0]
	fmt.Println(series.Uniform())
	// for _, c := series.Uniform() {
	// 	fmt.Println(c.WholesomeInfo())
	// }
}

func TestYcomb(t *testing.T) {
	a := assert.New(t)

	n1, err := html.Parse(strings.NewReader(`
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
`))
	a.NoError(err)

	n2, err := html.Parse(strings.NewReader(`
<html>
	<body>
		<div>
			<p>Hello 4</p>
		</div>
		<div>
			<p>Hello 5</p>
		</div>
		<div>
			<p>Hello 6</p>
		</div>
	</body>
</html>
`))
	a.NoError(err)

	arena := arena.NewArena()
	arena.Append(*n1)
	arena.Append(*n2)

	//s, _ := arena.RenderString(0)
	//fmt.Println(s)

	for i, item := range arena.List {
		if item.Data == "div" {
			fmt.Println(i)
		}
	}

	series := ExtractOptimized(arena).Matrix[0]

	for _, c := range series.Group.Clusters {
		v := c.TemplateVolume()
		if v > 0 {
			fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			fmt.Println(c.TemplateArena.RenderString(0))
			fmt.Println(":::::::::::::::::::::::::::::::::::::::::::::::::::::::")
			for _, ss := range c.Table {
				fmt.Println(ss)
			}
			fmt.Printf("similarity rate: %v, template volume: %v\n", c.Rate, v)
			fmt.Println("___________________________________________________")
		}
	}

	template := series.Nonuniform().Patterns()
	fmt.Printf("total chains: %v\n", len(template.Chains))
	for _, r := range template.Chains {
		fmt.Println(r.XPath())
	}
	fmt.Printf("size: %v, volume: %v, group volume: %v\n", series.Group.Size, series.Group.Volume, series.Group.GroupVolume)
	return
}

func TestOptimizedMatrix(t *testing.T) {
	a := assert.New(t)

	// f, err := os.Open("../fox.html")
	// a.NoError(err)
	// defer f.Close()
	// reader := bufio.NewReader(f)
	// n, err := html.Parse(reader)
	// a.NoError(err)

	n, err := html.Parse(strings.NewReader(`
		<html>
			<body>
				<div>
					<p>Hello 4</p>
				</div>
				<div>
					<p>Hello 5</p>
				</div>
				<div>
					<p>Hello 6</p>
				</div>
			</body>
		</html>
	`))
	a.NoError(err)

	arena := arena.NewArena()
	arena.Append(*n)
	Init(arena)

	s := NewDefaultComparator(arena)

	matrix1 := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float32 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	matrix2 := NewOptimizedRateMatrixAsync(len(arena.List), len(arena.List), 4, func(i, j int) float32 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	if len(matrix1.Values) != len(matrix2.Values) {
		t.Error("lengths differ")
	}

	for i := 0; i < len(matrix1.Values); i++ {
		for j := range matrix1.Values[i] {
			if matrix1.Get(i, j) != matrix2.Get(i, j) {
				t.Errorf("values differ: [%v][%v], %v and %v", i, j, matrix1.Get(i, j), matrix2.Get(i, j))
			}
		}
	}

	maxRate1, maxi1, maxj1 := matrix1.Max()
	maxRate2, maxi2, maxj2 := matrix2.Max()
	if maxRate1 != maxRate2 {
		t.Error("max rate differs")
	}
	if maxi1 != maxi2 {
		t.Error("maxi index differs")
	}
	if maxj1 != maxj2 {
		t.Error("maxj index differs")
	}
}

func TestOptimizedMatrixCloning(t *testing.T) {
	a := assert.New(t)

	// f, err := os.Open("../fox.html")
	// a.NoError(err)
	// defer f.Close()
	// reader := bufio.NewReader(f)
	// n, err := html.Parse(reader)
	// a.NoError(err)

	n, err := html.Parse(strings.NewReader(`
		<html>
			<body>
				<div>
					<p>Hello 4</p>
				</div>
				<div>
					<p>Hello 5</p>
				</div>
				<div>
					<p>Hello 6</p>
				</div>
			</body>
		</html>
	`))
	a.NoError(err)

	arena := arena.NewArena()
	arena.Append(*n)
	Init(arena)

	s := NewDefaultComparator(arena)
	matrix1 := NewRateMatrix(len(arena.List), len(arena.List), func(i, j int) float32 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})
	matrix2 := NewOptimizedRateMatrixAsync(len(arena.List), len(arena.List), 4, func(i, j int) float32 {
		if j <= i {
			return 0
		}
		return s.Cmp(s.arena.List[i], s.arena.List[j])
	})

	if len(matrix1.Values) != len(matrix2.Values) {
		t.Error("lengths differ")
	}

	for i := 0; i < len(matrix1.Values); i++ {
		for j := range matrix1.Values[i] {
			if matrix1.Get(i, j) != matrix2.Get(i, j) {
				t.Errorf("values differ: [%v][%v], %v and %v", i, j, matrix1.Get(i, j), matrix2.Get(i, j))
			}
		}
	}

	maxRate1, maxi1, maxj1 := matrix1.Max()
	maxRate2, maxi2, maxj2 := matrix2.Max()
	if maxRate1 != maxRate2 {
		t.Error("max rate differs")
	}
	if maxi1 != maxi2 {
		t.Error("maxi index differs")
	}
	if maxj1 != maxj2 {
		t.Error("maxj index differs")
	}

	for i := range matrix1.Values {
		idxs1 := matrix1.Candidates(i)
		idxs2 := matrix2.Candidates(i)
		if len(idxs1) != len(idxs2) {
			t.Error("indexes lengths differ")
		}
		for j := range idxs1 {
			if idxs1[j] != idxs2[j] {
				t.Errorf("index [%v][%v] differs", i, idxs1[j])
			}
		}
	}
}
