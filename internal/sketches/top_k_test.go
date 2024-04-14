package sketches

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

// Reuse tests from https://github.com/segmentio/topk/tree/main
func display(items []TopKItem) string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	for i, item := range items {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%s=%d", item.Key, item.Frequency)
	}
	buf.WriteRune(']')
	return buf.String()
}

func assert(t *testing.T, expect, given []TopKItem) {
	t.Helper()
	if len(expect) != len(given) {
		t.Fatalf("expected %d items, got %d:\nexpect: %s\ngot:    %s", len(expect), len(given), display(expect), display(given))
	} else {
		for i := range expect {
			if expect[i] != given[i] {
				t.Fatalf("element at %d doesn't match:\nexpect: %s\ngot:    %s", i, display(expect), display(given))
			}
		}
	}
}

func TestTopK(t *testing.T) {
	tests := []struct {
		desc   string
		k      int
		given  []TopKItem
		expect []TopKItem
	}{
		{
			desc:   "zero",
			k:      5,
			expect: []TopKItem{},
		},
		{
			desc: "simple, cardinality < k",
			k:    5,
			given: []TopKItem{
				{"c", 1},
				{"b", 5},
				{"a", 10},
				{"d", 25},
			},
			expect: []TopKItem{
				{"d", 25},
				{"a", 10},
				{"b", 5},
				{"c", 1},
			},
		},
		{
			desc: "simple, cardinality > k",
			k:    5,
			given: []TopKItem{
				{"c", 1},
				{"b", 5},
				{"a", 10},
				{"d", 25},
				{"f", 3},
				{"g", 20},
				{"h", 100},
				{"i", 2},
			},
			expect: []TopKItem{
				{"h", 100},
				{"d", 25},
				{"g", 20},
				{"a", 10},
				{"b", 5},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			topk := NewTopK(test.k)
			for _, fc := range test.given {
				for _ = range fc.Frequency {
					topk.Add(fc.Key)
				}
			}
			assert(t, test.expect, topk.List())
		})
	}
}

func TestSample_returnValue(t *testing.T) {
	topk := NewTopK(2)

	assert := func(key string, expect string) {
		t.Helper()
		got := topk.Add(key)
		if got != expect {
			t.Fatalf("for key %s, expected %v, got %v", key, expect, got)
		}
	}

	assert("a", "")
	assert("a", "")
	assert("a", "")

	assert("b", "")
	assert("b", "")

	assert("c", "")
	assert("c", "b")
}

func TestReset(t *testing.T) {
	topk := NewTopK(5)
	topk.Add("a")
	topk.Add("a")
	topk.Add("a")
	topk.Add("d")
	topk.Add("ad")
	topk.Reset()
	assert(t, []TopKItem{}, topk.List())
}

func BenchmarkSample(b *testing.B) {
	items := make([]string, 1_000_000)
	for i := range items {
		items[i] = randString(24)
	}

	for _, k := range []int{10, 50, 100, 500, 1_000, 5_000, 10_000, 1_000_000} {
		b.Run(fmt.Sprintf("K=%d", k), func(b *testing.B) {
			if len(items) < b.N {
				for i := len(items); i <= b.N; i++ {
					items = append(items, randString(16))
				}
			}
			items := make([]string, b.N)
			topk := NewTopK(k)
			b.ResetTimer()
			for _, key := range items[:b.N] {
				topk.Add(key)
			}
		})
	}
}

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}
