package bloblru

import (
	"testing"

	"github.com/restic/restic/internal/restic"
	rtest "github.com/restic/restic/internal/test"
)

func TestCache(t *testing.T) {
	var id1, id2, id3 restic.ID
	id1[0] = 1
	id2[0] = 2
	id3[0] = 3

	const (
		kiB       = 1 << 10
		cacheSize = 64*kiB + 3*overhead
	)

	c := New(cacheSize)

	addAndCheck := func(id restic.ID, exp []byte) {
		c.Add(id, exp)
		blob, ok := c.Get(id)
		rtest.Assert(t, ok, "blob %v added but not found in cache", id)
		rtest.Equals(t, &exp[0], &blob[0])
		rtest.Equals(t, exp, blob)
	}

	addAndCheck(id1, make([]byte, 32*kiB))
	addAndCheck(id2, make([]byte, 30*kiB))
	addAndCheck(id3, make([]byte, 10*kiB))

	_, ok := c.Get(id2)
	rtest.Assert(t, ok, "blob %v not present", id2)
	_, ok = c.Get(id1)
	rtest.Assert(t, !ok, "blob %v present, but should have been evicted", id1)

	c.Add(id1, make([]byte, 1+c.size))
	_, ok = c.Get(id1)
	rtest.Assert(t, !ok, "blob %v too large but still added to cache")

	c.c.Remove(id1)
	c.c.Remove(id3)
	c.c.Remove(id2)

	rtest.Equals(t, cacheSize, c.size)
	rtest.Equals(t, cacheSize, c.free)
}
