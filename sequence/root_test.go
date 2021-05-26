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

	r.Batch().Results()

	for _, c := range r.clusters {
		if len(c.indexes) != len(c.stemIndexes) {
			t.Error("bad")
		}
	}
}

func TestRootCluster_Batch2(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadString(testDoc2)
	a.NoError(err)

	for _, s := range r.Batch().Results() {
		fmt.Println(s.TransposedValues)
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
			fmt.Println(cc.indexes)
			fmt.Println(cc.rate)
		}
	}
	fmt.Println()
	for _, c := range r2.clusters {
		for _, cc := range c.clusters {
			fmt.Println(cc.indexes)
			fmt.Println(cc.rate)
		}
	}
}

func TestRootCluster_LoadFile(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster() //.SetLimit(10)
	err := r.LoadFile("./rozetka1.html")
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
	r := NewRootCluster() //.SetLimit(10)
	//err := r.LoadFile("../bin/samples/ycomb.html")
	//err := r.LoadFile("../bin/samples/hn.html")
	err := r.LoadFile("../bin/samples/pravda1.html")
	a.NoError(err)
	series := r.Batch().Results()
	for i, s := range series {
		if s.Group.Size == 40 {
			fmt.Println(i)
		}
	}
}

func TestRootCluster_LoadMultipleFiles(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster() //.SetLimit(20)
	err := r.LoadFile("./rozetka1.html")
	a.NoError(err)
	err = r.LoadFile("./rozetka2.html")
	a.NoError(err)
	err = r.LoadFile("./rozetka1.html")
	a.NoError(err)

	series := r.Batch().Results()
	for i, s := range series {
		if len(s.TransposedValues) == 180 {
			fmt.Println(i)
		}
	}
}

func TestRootCluster_FB(t *testing.T) {
	a := assert.New(t)

	r := NewRootCluster()
	err := r.LoadFile("./fb.html")
	a.NoError(err)

	fmt.Println(r.Arena.Find("div", "data-pagelet", "FeedUnit_{n}"))

	r.Batch().Results()

	for _, c := range r.nodeIDToCluster[2694].clusters {
		if c.Has(2694) {
			fmt.Println(c.Volume())
			fmt.Println(20.41* float64(len(c.indexes)-5))
			for i := range c.indexes[1:] {
				//fmt.Printf("%v:%v \n", c.indexes[i], c.indexes[i+1])
				n, m := r.nodeIDToCluster[2694].FindIdx(c.indexes[i]), r.nodeIDToCluster[2694].FindIdx(c.indexes[i+1])
				fmt.Printf("%v:%v = %v\n", c.indexes[i], c.indexes[i+1], r.nodeIDToCluster[2694].Get(n, m))
			}
		}
	}

	//for _, s := range series {
	//	if s.Group.Clusters[0].Members[0].GetAttr("data-pagelet") != "" {
	//		fmt.Println("got it")
	//	}
	//	//if s.Group.Size == 63 {
	//	//	for _, t := range s.Group.Clusters {
	//	//		fmt.Println(t.Members)
	//	//	}
	//	//}
	//}
}