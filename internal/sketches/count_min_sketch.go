package sketches

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// Count min sketch algorithm implementation. Using murmur hash for good distribution
// use as reference http://dimacs.rutgers.edu/~graham/pubs/papers/cm-full.pdf

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// CMS structure
type CountMinSketch struct {
	width  uint32
	depth  uint32
	counts [][]uint32
}

// CMS constructor
func NewCountMinSketch(width, depth uint32) *CountMinSketch {

	c := &CountMinSketch{}
	c.width = width
	c.depth = depth
	c.counts = make([][]uint32, depth)
	for col := range c.counts {
		c.counts[col] = make([]uint32, c.width)
	}
	return c
}

func (c *CountMinSketch) Add(item string) {
	for i := uint32(0); i < c.depth; i++ {
		bucketIndex := murmurHash(item, i) % uint32(c.width)
		c.counts[i][bucketIndex] += 1
	}
}

func (c *CountMinSketch) Count(item string) uint32 {
	res := uint32(math.MaxUint32)
	for i := uint32(0); i < c.depth; i++ {
		bucketIndex := murmurHash(item, i) % c.width
		res = min(res, c.counts[i][bucketIndex])
	}
	return res
}

func (c *CountMinSketch) Reset() {
	for col := range c.counts {
		for i := range c.counts[col] {
			c.counts[col][i] = 0
		}
	}
}

// murmurHash generates a hash value for the given item and seed using MurmurHash3
func murmurHash(item string, seed uint32) uint32 {
	hash := murmur3.New32WithSeed(seed)
	hash.Write([]byte(item))
	return hash.Sum32()
}
