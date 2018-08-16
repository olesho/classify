// go-words project words.go
package words

import (
	"fmt"
	"sort"
)

type Char byte

func (c1 Char) SameKind(c2 Comparable) bool {
	if _, ok := c2.(Char); ok {
		return c1 == c2
	}
	panic("Comparing incompatible interfaces")
	return false
}

type Str string

func (c1 Str) SameKind(c2 Comparable) bool {
	if _, ok := c2.(Str); ok {
		return c1 == c2
	}
	panic("Comparing incompatible interfaces")
	return false
}

type Comparable interface {
	SameKind(Comparable) bool
}

func (w1 Word) SameKind(w2 Word) bool {
	if len(w1) == len(w2) {
		for i, _ := range w1 {
			if !w1[i].SameKind(w2[i]) {
				return false
			}
		}
		return true
	}
	return false
}

type Word []Comparable
type Positions []int
type Tuple struct {
	Word
	Positions
}

func (t *Tuple) String() string {
	return fmt.Sprint(t.Word) + fmt.Sprint(t.Positions)
}

type Processor struct {
	Text       Word
	vocabulary []Tuple
	current    []Tuple
}

func NewProcessor() *Processor {
	return &Processor{
		Text:       Word{},
		vocabulary: make([]Tuple, 0),
		current:    make([]Tuple, 0),
	}
}

func (p *Processor) Done() {
	for i := 0; i < len(p.current); i++ {
		p.saveTuple(p.current[i])
	}
	p.current = []Tuple{}
}

func (p *Processor) Next(nextItem Comparable) {
	p.Text = append(p.Text, nextItem)
	index := p.wordIndexInVocabulary(nextItem)

	p.saveTuple(Tuple{Word{nextItem}, []int{len(p.Text) - 1}})
	if index > -1 {
		toExclude := make([]int, 0)
		for i := 0; i < len(p.current); i++ {
			pos := p.extends(p.current[i], nextItem)
			if len(pos) < 2 {
				p.current[i].Positions = append(p.current[i].Positions, len(p.Text)-len(p.current[i].Word))
				p.saveTuple(p.current[i])
				p.current = append(p.current[:i], p.current[i+1:]...)
				i--
			} else {
				toExclude = appendUnique(toExclude, pos)

				p.current[i].Word = append(p.current[i].Word, nextItem)
				p.current[i].Positions = pos

				if p.current[i].Word.repeats() {
					p.saveTuple(Tuple{p.current[i].Word.half(), p.current[i].Positions}) //unique(p.current[0].Positions, p.current[1].Positions)})
					p.current = append(p.current[:i], p.current[i+1:]...)
					i--
				}
			}
		}

		// append current character as current word
		t := p.vocabulary[index]
		if (len(t.Positions) - len(toExclude)) > 0 {
			p.current = append(p.current, t)
		}
		return
	}

	for _, t := range p.current {
		p.saveTuple(t)
	}
	p.current = []Tuple{Tuple{Word{nextItem}, []int{len(p.Text) - 1}}}
}

func (p *Processor) includes(t2 Tuple) bool {
	if len(p.current) > 0 {
		for _, t1 := range p.current {
			if t1.includes(t2) {
				return true
			}
		}
	}
	return false
}

func (t1 *Tuple) includes(t2 Tuple) bool {
	if len(t1.Positions) == len(t2.Positions) {
		offset := t2.Positions[0] - t1.Positions[0]
		for i, _ := range t1.Positions {
			if t1.Positions[i]+offset != t2.Positions[i] {
				return false
			}
		}
		return true
	}
	return false
}

func (p *Processor) saveTuple(t Tuple) {
	for i, currentTuple := range p.vocabulary {
		if currentTuple.Word.SameKind(t.Word) {
			p.vocabulary[i].Positions = unique(t.Positions, p.vocabulary[i].Positions)
			return
		}
	}
	p.vocabulary = append(p.vocabulary, t)
}

func (p *Processor) extends(t Tuple, nextItem Comparable) []int {
	r := []int{}
	if len(t.Positions) > 1 {
		pos := t.Positions[0] + len(t.Word)
		if p.Text[pos].SameKind(nextItem) {
			r = append(r, t.Positions[0])
		}

		for i := 1; i < len(t.Positions); i++ {
			pos := t.Positions[i] + len(t.Word)
			if p.Text[pos].SameKind(nextItem) {
				r = append(r, t.Positions[i])
			}
		}
	}
	return r
}

func (w1 Word) repeats() bool {
	if len(w1)%2 == 0 {
		if Word(w1[:len(w1)/2]).SameKind(Word(w1[len(w1)/2:])) {
			return true
		}
	}
	return false
}

func (w1 Word) half() Word {
	return Word(w1[:len(w1)/2])
}

func substract(from, pos []int) []int {
	res := make([]int, 0)
	for _, f := range from {
		if !in(f, pos) {
			res = append(res, f)
		}
	}
	return res
}

func unique(p1, p2 []int) []int {
	p3 := make([]int, len(p1), len(p1)+len(p2))
	copy(p3, p1)
	for _, p := range p2 {
		if !in(p, p3) {
			p3 = append(p3, p)
		}
	}
	return p3
}

func appendUnique(dest, src []int) []int {
	for _, s := range src {
		if !in(s, dest) {
			dest = append(dest, s)
		}
	}
	return dest
}

func in(n int, pos []int) bool {
	for _, p := range pos {
		if p == n {
			return true
		}
	}
	return false
}

func (p *Processor) wordIndexInVocabulary(ch Comparable) int {
	for i, wr := range p.vocabulary {
		if len(wr.Word) == 1 {
			if wr.Word[0].SameKind(ch) {
				return i
			}
		}
	}
	return -1
}

func (p *Processor) Vocabulary() []Tuple {
	return p.vocabulary
}

func (p *Processor) SortVocabulary() {
	sort.Slice(p.vocabulary, func (i, j int) bool {
		return len(p.vocabulary[i].Word)*p.Count(Word(p.vocabulary[i].Word)) > len(p.vocabulary[j].Word)*p.Count(Word(p.vocabulary[j].Word))
	})
}


func (p *Processor) Count(word Word) int {
	cnt := 0
	for i := 0; i < len(p.Text); i++ {
		if startsWith(Word(p.Text[i:]), word) {
			cnt++
			i += len(word)-1
		}
	}
	return cnt
}

func FindPositions(word, text Word) []int {
	res := []int{}
	for i := 0; i < len(text); i++ {
		if startsWith(Word(text[i:]), word) {
			res = append(res, i)
		}
	}
	return res
}

func startsWith(w1, w2 Word) bool {
	if len(w2) > len(w1) {
		return false
	}
	for i, _ := range w2 {
		if !w1[i].SameKind(w2[i]) {
			return false
		}
	}
	return true
}