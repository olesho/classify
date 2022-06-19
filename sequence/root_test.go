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
				<div>127.0.0.1</div>
				<h1>some@gmail.com</h1>
            </div>
            <div>
                <h3> Title 2 </h3>
                <p> Some text 2 </p>
                <img src="/src2"> image2 </img>
				<div>127.0.0.2</div>
				<h1>one@gmail.com</h1>
            </div>
            <div>
                <h3> Title 3 </h3>
                <p> Some text 3 </p>
                <img src="/src3"> image3 </img>
				<div>127.0.0.3</div>
				<h1>two@gmail.com</h1>
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
				<div>127.0.0.4</div>
				<h1>three@gmail.com</h1>
            </div>
            <div>
                <h3> Title 5 </h3>
                <p> Some text 5 </p>
                <img src="/src5"> image2 </img>
				<div>127.0.0.5</div>
				<h1>four@gmail.com</h1>
            </div>
            <div>
                <h3> Title 6 </h3>
                <p> Some text 6 </p>
                <img src="/src6"> image3 </img>
				<div>127.0.0.6</div>
				<h1>five@gmail.com</h1>
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
		types := s.GetFieldTypes()
		fmt.Println(types)
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

	for _, c := range r2.clusters {
		for _, cc := range c.clusters {
			fmt.Println(cc.items)
			fmt.Println(cc.rate)
		}
	}
}

func Test_SmallFile(t *testing.T) {
	fileName := "ips_small.html"
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadFile(fileName)
	a.NoError(err)
	series := r.Batch().Results()
	trFound := r.Arena.FindByName("tr")
	for _, s := range series {
		if len(s.TransposedValues) == len(trFound) {
			return
		}
	}
	t.Error("no success")
}

func TestRootCluster_LoadTinyFile(t *testing.T) {
	// change to file input
	fileName := "ips_tiny.html"
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadFile(fileName)
	a.NoError(err)
	series := r.BatchSync().Results()
	trFound := r.Arena.FindByName("tr")
	for _, s := range series {
		if len(s.TransposedValues) == len(trFound) {
			return
		}
	}
	t.Error("no success")
}

func TestRootCluster_LoadAvgFile(t *testing.T) {
	// change to file input
	fileName := "ips_avg.html"
	a := assert.New(t)
	r := NewRootCluster().SetDebug(&DebugConfig{TagName: "tr"})
	err := r.LoadFile(fileName)
	a.NoError(err)
	series := r.Batch().Results()
	trFound := r.Arena.FindByName("tr")
	for _, s := range series {
		if len(s.TransposedValues) == len(trFound) {
			return
		}
	}
	t.Error("no success")
}

func TestRootCluster_LoadPravda(t *testing.T) {
	fileName := "pravda.html"
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadFile(fileName)
	a.NoError(err)
	series := r.Batch().Results()
	for _, s := range series {
		if len(s.TransposedValues) == 40 {
			return
		}
	}
	t.Error("no success")
}

func TestRootCluster_LoadHackernoon(t *testing.T) {
	fileName := "hackernoon.html"
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadFile(fileName)
	a.NoError(err)
	series := r.Batch().Results()
	for _, s := range series {
		if len(s.TransposedValues) == 63 {
			return
		}
	}
	t.Error("no success")
}

func TestRootCluster_LoadMini(t *testing.T) {
	a := assert.New(t)
	r := NewRootCluster()
	err := r.LoadString(`
		<html>
			<body>
				<p>
					<div>txt1</div>
					<h1>head1<h1>
					<div>txt2</div>
				</p>
				<p>
					<div>txt2</div>
					<h1>head2<h1>
					<div></div>
				</p>
				<p>
					<div></div>
					<h1>head3<h1>
					<div>txt3</div>
				</p>
			</body>
		</html>
`)
	a.NoError(err)
	series := r.BatchSync().Results()
	for _, s := range series {
		if len(s.TransposedValues) == 3 {
			return
		}
	}
	t.Error("no success")
}

//func TestChainsOptimization(t *testing.T) {
//	a := assert.New(t)
//	r := NewRootCluster()
//	err := r.LoadString(testDoc1)
//	a.NoError(err)
//
//	for i := range r.Arena.List {
//		chain :=  r.Arena.Chain(i, 0)
//		chainIDXs := r.Arena.ChainIDXs(i, 0)
//		l1 := len(chain)
//		l2 := len(chainIDXs)
//		if l1 != l2 {
//			t.Error("different length")
//			return
//		}
//
//		for i, c := range chain {
//			if c.Id != chainIDXs[i] {
//				t.Error("different indexes")
//				return
//			}
//		}
//	}
//}