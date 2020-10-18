package stream

import (
	"bufio"
	"github.com/olesho/classify/arena"
	"github.com/olesho/classify/comparator"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

type Finder interface {
	Find(idx1, idx2 int) float32
}

func (s *Storage) TryAdd(clusterIdx, idx int) bool {
	mtx := s.Clusters[clusterIdx]
	if s.StrictComparator.Cmp(mtx.Indexes[0], idx) > 0 {
		mtx.mutex.Lock()
		defer mtx.mutex.Unlock()

		values := make([]float32, len(mtx.Indexes))
		for i, existingIdx := range mtx.Indexes {
			values[i] = s.ElementComparator.Cmp(idx, existingIdx)
			if values[i] == 0 {
				return false
			}
		}
		for i := range mtx.Indexes {
			mtx.Values[i] = append(mtx.Values[i], values[i])
		}
		mtx.Indexes = append(mtx.Indexes, idx)
		mtx.Values = append(mtx.Values, []float32{})
		s.NodeToCluster[idx] = mtx
		return true
	}
	return false
}

type Storage struct {
	Arena *arena.Arena
	Clusters []*Mtx
	NodeToCluster []*Mtx
	StrictComparator comparator.Comparator
	ElementComparator comparator.Comparator
	//mutex sync.Mutex
	timer *Timer
}

func NewStorage() *Storage {
	a := arena.NewArena()
	return &Storage{
		Arena: a,
		StrictComparator: comparator.NewStrictComparator(a),
		ElementComparator: comparator.NewElementComparator(a),
		timer: NewTimer(),
	}
}

// Load appends HTML file content from reader
func (s *Storage) Load(r io.Reader) error {
	n, err := html.Parse(r)
	if err != nil {
		return err
	}
	s.Arena.Append(*n)
	return nil
}

// LoadFile appends HTML file content
func (s *Storage) LoadFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	n, err := html.Parse(reader)
	if err != nil {
		return err
	}
	s.Arena.Append(*n)
	return nil
}

// LoadString appends HTML string
func (s *Storage) LoadString(str string) error {
	n, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return err
	}
	s.Arena.Append(*n)
	return nil
}

func (s *Storage) next(idx int) {
	var i int
	for i = 0; i < len(s.Clusters) && !s.TryAdd(i, idx); i++ {}
	if i == len(s.Clusters) {
		s.Clusters = append(s.Clusters, s.NewMtx(idx))
	}
}

func (s *Storage) clusterRun(clusterIdx int) {
	mtx := s.Clusters[clusterIdx]
	for i := range mtx.Values {
		idx1 := mtx.Indexes[i]
		for j, idx2 := range mtx.Indexes[i+1:] {
			mtx.Values[i][j] += s.cmpChildren(idx1, idx2)
		}
	}
}


func (s *Storage) Find(idx1, idx2 int) float32 {
	c1, c2 := s.NodeToCluster[idx1], s.NodeToCluster[idx2]
	if c1 != c2 {
		return 0
	}

	i, j := -1, -1
	for n, idx := range c1.Indexes {
		if idx == idx1 {
			i = n
			break
		}
	}
	for n, idx := range c1.Indexes {
		if idx == idx2 {
			j = n
			break
		}
	}
	return c1.Get(i, j)
}