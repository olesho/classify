package sequence

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDoc1 = `
	<html>
		<body>
			<section>
				<div>
					<p>Hello 1</p>
				</div>
				<div>
					<p>Hello 2</p>
				</div>
				<div>
					<p>Hello 3</p>
				</div>
			</section>
			<p> some garbage ad </p>
			<span> some garbage ad </span>
			<b>  some garbage ad </b>
			<div>  some garbage ad </div>
			<h1>  some garbage ad  </h1>
			<h2>  some garbage ad  </h2>
			<h3>  some garbage ad  </h3>
			<section>
				<div>
					<p>Hello 4</p>
				</div>
				<div>
					<p>Hello 5</p>
				</div>
				<div>
					<p>Hello 6</p>
				</div>
			</section>
			<p> some garbage ad </p>
			<span> some garbage ad </span>
			<b>  some garbage ad </b>
			<div>  some garbage ad </div>
			<h1>  some garbage ad  </h1>
			<h2>  some garbage ad  </h2>
			<h3>  some garbage ad  </h3>
			<section>
				<div>
					<p>Hello 7</p>
				</div>
				<div>
					<p>Hello 9</p>
				</div>
				<div>
					<p>Hello 10</p>
				</div>
			</section>
		</body>
	</html>
`

func TestRootCluster_Batch(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadString(testDoc1)
	a.NoError(err)

	r.Batch().Results()

	for _, c := range r.clusters {
		if len(c.indexes) != len(c.stemIndexes) {
			t.Error("bad")
		}
	}

	fmt.Println(r)
}

func TestRootCluster_LoadFile(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()//.SetLimit(10)
	err := r.LoadFile("../bin/samples/rozetka.html")
	a.NoError(err)
	series := r.Batch().Results()
	for i, s := range series {
		if s.Group.Size == 60 {
			fmt.Println(i)
		}
	}

	// TODO: get XPath
	for _, n := range series[0].TransposedNodes {
		fmt.Println(n)
	}

}

func TestRootCluster_Hackernews(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()//.SetLimit(10)
	//err := r.LoadFile("../bin/samples/ycomb.html")
	err := r.LoadFile("../bin/samples/hn.html")
	a.NoError(err)
	series := r.Batch().Results()
	for i, s := range series {
		if s.Group.Size == 30 {
			fmt.Println(i)
		}
	}
}

func TestRootCluster_LoadMultipleFiles(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster().SetLimit(20)
	err := r.LoadFile("../rozetka1.html")
	a.NoError(err)
	err = r.LoadFile("../rozetka2.html")
	a.NoError(err)
	err = r.LoadFile("../rozetka3.html")
	a.NoError(err)

	series := r.Batch().Results()
	for i, s := range series {
		if s.Group.Size == 180 {
			fmt.Println(i)
		}
	}
}