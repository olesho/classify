package stream

import (
	"fmt"
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

func TestStorageShort(t *testing.T) {
	s := NewStorage()
	err := s.LoadString(testDoc1)
	if err != nil {
		t.Error(err)
	}
	s.Run()
	for _, c := range s.Clusters {
		fmt.Println(c.Indexes)
	}
}

func _equalMtx(m1, m2 *Mtx) bool {
	if len(m1.Indexes) != len(m2.Indexes) {
		return false
	}

	for i, idx1 := range m1.Indexes {
		if !m2.HasIdx(idx1) {
			return false
		}
		if len(m1.Values[i]) != len(m2.Values[i]) {
			return false
		}
	}
	return true
}

func _hasEqualMtx(list []*Mtx, m2 *Mtx) bool {
	for _, m1 := range list {
		if _equalMtx(m1, m2) {
			return true
		}
	}
	return false
}

func TestStorageSyncAsyncShort(t *testing.T) {
	syncS := NewStorage()
	err := syncS.LoadString(testDoc1)
	if err != nil {
		t.Error(err)
	}
	syncS.Run()

	asyncS := NewStorage()
	err = asyncS.LoadString(testDoc1)
	if err != nil {
		t.Error(err)
	}
	asyncS.RunAsync()

	if len(syncS.Clusters) != len(asyncS.Clusters) {
		t.Error("cluster count differ")
	}

	for _, c1 := range asyncS.Clusters {
		if !_hasEqualMtx(syncS.Clusters, c1) {
			t.Error("missing matrix")
		}
	}
}

func TestStorageSyncAsyncLong(t *testing.T) {
	syncS := NewStorage()
	err := syncS.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	syncS.Run()

	asyncS := NewStorage()
	err = asyncS.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	asyncS.RunAsync()

	if len(syncS.Clusters) != len(asyncS.Clusters) {
		t.Errorf("cluster count differ: %v vs %v", len(syncS.Clusters), len(asyncS.Clusters))
	}

	for _, c1 := range asyncS.Clusters {
		if !_hasEqualMtx(syncS.Clusters, c1) {
			t.Error("missing matrix")
		}
	}
}

func TestBasicStorageLong(t *testing.T) {
	s := NewStorage()
	err := s.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	s.RunAsync()
	//matrix := s.RunAsync()
	//for _, series := range matrix.Matrix {
	//	fmt.Println(series.Group.Size)
	//}
	//
	//for i, c := range s.Clusters {
	//	if c.HasIdx(5715) {
	//		fmt.Println(i)
	//	}
	//
	//	//if len(c.indexes) == 60 {
	//	//	fmt.Println(s.Arena.Get(c.indexes[0]).String())
	//	//}
	//}
}

func TestMtx(t *testing.T) {
	mtx := &Mtx{}
	err := mtx.load("div.mtx")
	if err != nil {
		t.Error(err)
	}

	s := NewStorage()
	err = s.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}

	clusters := mtx.GenerateClusters()
	for _, c := range clusters {
		idx := c.Indexes[0]
		classes := s.Arena.Get(idx).Classes()
		class := ""
		if len(classes) > 0 {
			class = classes[0]
		}
		if class == "catalog-grid__cell" {
			fmt.Println(len(c.Indexes), c.Rate)
		}

		//el := make([]string, len(c.Indexes))
		//for i, idx := range c.Indexes {
		//	classes := s.Arena.Get(idx).Classes()
		//	class := ""
		//	if len(classes) > 0 {
		//		class = classes[0]
		//	}
		//	el[i] = fmt.Sprintf("%v-%v", idx, class)
		//	if class == "catalog-grid__cell" {
		//		fmt.Println(len(c.Indexes), c.Rate)
		//	}
		//}
		//fmt.Println(el)
	}
}