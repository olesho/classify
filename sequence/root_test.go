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

var testDoc2 = `
<html>
    <body>
        <section> Some Ad </section>
        <section> 
            <h1> Data </h1> 
            <div>
                <h3> Title 1 </h3>
                <p> Some text 1 </p>
                <img src="/src1"> image1 </img>
            </div>
            <div>
                <h3> Title 2 </h3>
                <p> Some text 2 </p>
                <img src="/src2"> image2 </img>
            </div>
            <div>
                <h3> Title 3 </h3>
                <p> Some text 3 </p>
                <img src="/src3"> image3 </img>
            </div>
        </section>
        <section> 
            <h2> Some Menu </h2>
            <ul>
                <li>Item 1</li>
                <li>Item 2</li>
                <li>Item 3</li>
            </ul>
        </section>
    </body>
</html>
`

var testDoc3 = `
<html>
    <body>
        <section> Some Ad </section>
        <section> 
            <h1> Data </h1> 
            <div>
                <h3> Title 4 </h3>
                <p> Some text 4 </p>
                <img src="/src4"> image1 </img>
            </div>
            <div>
                <h3> Title 5 </h3>
                <p> Some text 5 </p>
                <img src="/src5"> image2 </img>
            </div>
            <div>
                <h3> Title 6 </h3>
                <p> Some text 6 </p>
                <img src="/src6"> image3 </img>
            </div>
        </section>
        <section> 
            <h2> Some Menu </h2>
            <ul>
                <li>Item 1</li>
                <li>Item 2</li>
                <li>Item 3</li>
            </ul>
        </section>
    </body>
</html>
`

func TestRootCluster_Batch(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadString(testDoc1)
	a.NoError(err)

	//r.Batch().Results()
	r.BatchSync().Results()

	for _, c := range r.clusters {
		if len(c.indexes) != len(c.stemIndexes) {
			t.Error("bad")
		}
	}
}


func TestRootCluster_BatchMultiple(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadString(testDoc2)
	a.NoError(err)

	err = r.LoadString(testDoc3)
	a.NoError(err)

	for _, s := range r.Batch().Results() {
		fmt.Println(s.TransposedValues)
	}
}

func TestRootCluster_Batch2SyncAsync(t *testing.T) {
	a := assert.New(t)

	r1 := NewRootCluster()
	err := r1.LoadString(testDoc2)
	a.NoError(err)
	r1.BatchSync()

	r2 := NewRootCluster()
	err = r2.LoadString(testDoc2)
	a.NoError(err)
	r2.Batch()

	for _, c := range r1.clusters {
		for _, cc := range c.clusters {
			fmt.Println(cc.items)
			fmt.Println(cc.rate)
		}
	}
	fmt.Println()
	for _, c := range r2.clusters {
		for _, cc := range c.clusters {
			fmt.Println(cc.items)
			fmt.Println(cc.rate)
		}
	}
}

func TestRootCluster_LoadFile(t *testing.T) {
	// change to file input
	fileName := ""

	if fileName != "" {
		a := assert.New(t)
		r := NewRootCluster() //.SetLimit(10)
		err := r.LoadFile(fileName)
		a.NoError(err)
		series := r.Batch().Results()
		for i, s := range series {
			fmt.Println(i, s.TransposedValues)
		}
	}
}
