package sequence

import "strings"

type CrownCluster struct {
	indexes []int
	stem *StemCluster
	rate float32
}

func (c *CrownCluster) String() string {
	r := make([]string, len(c.indexes))
	for i, index := range c.indexes {
		r[i] = c.stem.arena.StringifyNode(index)
	}
	return strings.Join(r, " ")
}

func (c *CrownCluster) Volume() float32 {
	return float32(len(c.indexes)) * c.rate
}

func (c *CrownCluster) Add(rate float32, index int) bool {
	if c.Volume() < rate * float32(len(c.indexes)+1) {
		c.rate = rate
		c.indexes = append(c.indexes, index)
		return true
	}
	return false
}

func (c *CrownCluster) Rate(stemIndex int) float32 {
	var lowestVal float32
	for i := range c.indexes {
		v := c.stem.Get(stemIndex, i)
		if v > 0 {
			if lowestVal == 0 {
				lowestVal = v
				continue
			}
			if v < lowestVal {
				lowestVal = v
			}
		} else {
			return 0
		}
	}
	return lowestVal
}