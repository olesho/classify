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


//
//func TestStorageShort(t *testing.T) {
//	s := NewStorage()
//	err := s.LoadString(testDoc1)
//	if err != nil {
//		t.Error(err)
//	}
//	s.Run()
//	for _, c := range s.Clusters {
//		fmt.Println(c.Indexes)
//	}
//}

func _equalMtx(m1, m2 *Mtx) bool {
	if len(m1.Indexes) != len(m2.Indexes) {
		return false
	}

	for i, idx1 := range m1.Indexes {
		if m2.FindIdx(idx1) == -1 {
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
//
//func TestStorageSyncAsyncShort(t *testing.T) {
//	syncS := NewStorage()
//	err := syncS.LoadString(testDoc1)
//	if err != nil {
//		t.Error(err)
//	}
//	syncS.Run()
//
//	asyncS := NewStorage()
//	err = asyncS.LoadString(testDoc1)
//	if err != nil {
//		t.Error(err)
//	}
//	asyncS.RunAsync()
//
//	if len(syncS.Clusters) != len(asyncS.Clusters) {
//		t.Error("cluster count differ")
//	}
//
//	for _, c1 := range asyncS.Clusters {
//		if !_hasEqualMtx(syncS.Clusters, c1) {
//			t.Error("missing matrix")
//		}
//	}
//}
//
func TestStorageSync(t *testing.T) {
	syncS := NewStorage()
	err := syncS.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	syncS.createMatrices()
	syncS.compareInMatrices()
	if syncS.NodeToCluster[5715] == syncS.NodeToCluster[6627] {
		syncS.NodeToCluster[5715].save("li.mtx")
	} else {
		t.Error("clusters should be equal for selected node")
	}
}

func TestCreateMatricesSyncAsync(t *testing.T) {
	syncS := NewStorage()
	err := syncS.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	syncS.createMatrices()

	asyncS := NewStorage()
	err = asyncS.LoadFile("../rozetka.html")
	if err != nil {
		t.Error(err)
	}
	asyncS.createMatrices()

	for i := range syncS.Arena.List {
		if !syncS.NodeToCluster[i].Equal(asyncS.NodeToCluster[i]) {
			t.Errorf("matrix %v not equals sync to async", i)
		}
	}

	syncS.compareInMatrices()
	asyncS.compareInMatricesAsync()
}

func TestMtx(t *testing.T) {
	s := NewStorage()
	err := s.LoadFile("./testDoc2.html")
	if err != nil {
		t.Error(err)
	}
	s.timer.Start()
	s.createMatricesAsync()
	s.compareInMatricesAsync()
	clusters := s.generateAllClustersAsync()
	for _, c := range clusters {
		if len(c.Indexes) == 60 {
			table := s.toTable(c)
			fmt.Println(table)
		}
	}

	series := s.clustersToMatrix(clusters)

	for _, item := range series {
		if item.Group.Size == 60 {
			for _, cluster := range item.Group.Clusters {
				for _, field := range cluster.Fields {
					fmt.Println(field)
				}
			}
		}
	}

	//for _, mtx := range s.Clusters {
	//	if mtx.FindIdx(s.Arena.FindByAttr("class", "catalog-grid__cell  catalog-grid__cell_type_slim")[0]) > -1 {
	//		//clusters := mtx._GenerateClusters(s.Arena)
	//		clusters := mtx._GenerateClusters(s.Arena)
	//		filtered := []*Cluster{}
	//		for _, c := range clusters {
	//			//fmt.Println(s.Arena.List[c.Indexes[0]].Classes())
	//			if s.Arena.List[c.Indexes[0]].HasClass("catalog-grid__cell_type_slim") {
	//				fmt.Println("before", len(c.Indexes))
	//				filtered = append(filtered, c)
	//			}
	//		}
	//
	//		clusters = s.consolidateAll(filtered)
	//		for _, c := range clusters {
	//			//fmt.Println(s.Arena.List[c.Indexes[0]].Classes())
	//			if s.Arena.List[c.Indexes[0]].HasClass("catalog-grid__cell_type_slim") {
	//				fmt.Println("after", len(c.Indexes))
	//			}
	//		}
	//	}
	//}
}

func TestSyncAsync1(t *testing.T) {
	s1 := NewStorage()
	err := s1.LoadFile("./testDoc1.html")
	if err != nil {
		t.Error(err)
	}
	s1.createMatricesAsync()
	s1.compareInMatricesAsync()

	s2 := NewStorage()
	err = s2.LoadFile("./testDoc1.html")
	if err != nil {
		t.Error(err)
	}
	s2.createMatrices()
	s2.compareInMatrices()

	if s2.Find(62, 150) != s1.Find(62, 150) {
		t.Error("Find(62, 150) error")
	}
	if s2.Find(63, 151) != s1.Find(63, 151) {
		t.Error("Find(63, 151)")
	}
	if s2.Find(72, 151) != s1.Find(72, 151) {
		t.Error("Find(72, 151)")
	}
}

func TestSyncAsync2(t *testing.T) {
	s1 := NewStorage()
	err := s1.LoadFile("./testDoc2.html")
	if err != nil {
		t.Error(err)
	}
	s1.timer.Start()
	s1.createMatricesAsync()
	s1.compareInMatricesAsync()
	clusters := s1.generateAllClustersAsync()
	ids := s1.Arena.FindByAttr("class", "catalog-grid__cell  catalog-grid__cell_type_slim")
	for _, c := range clusters {
		if c.hasIndex(ids[0]) {
			if len(c.Indexes) != 60 {
				t.Error("wrong group generated")
			}
		}
	}
}
