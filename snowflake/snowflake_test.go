package snowflake

import (
	"math/rand"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	i0 := NextID()
	doTest(t, "NextID", i0)
	i0 = Parse(NextInt(), 0)
	doTest(t, "NextInt", i0)
	i0 = ParseFormat(NextFormat(0), 0, 0)
	doTest(t, "NextFormat", i0)
	SetDefault(func(node *Snowflake) {
		node.SetEpoch(time.Unix(0, 0))
		node.SetNode(3)
	})
	i0 = ParseString(NextString(), 0)
	doTest(t, "NextString", i0)

	i0 = ParseString(New().Next().String(), 0)
	doTest(t, "NewNextString", i0)

	//测序
	n := New()

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		var preId ID
		idx := rnd.Int63n(1000000000)
		for j := idx; j < idx+10000; j++ {
			curId := n.Next()
			if j > idx {
				if curId.Int() < preId.Int() || curId.String() < preId.String() || curId.Format(0) < preId.Format(0) {
					t.Logf("pre -> %09d %d %s %s", j-1, preId.Int(), preId.String(), preId.Format(0))
					t.Logf("cur -> %09d %d %s %s", j, curId.Int(), curId.String(), curId.Format(0))
					t.Fail()
				}
			}
			preId = curId
		}
	}
}

func doTest(t *testing.T, s string, i0 ID) {
	i1 := ParseFormat(Parse(i0.Int(), 0).Format(0), 0, 0)
	Tes(t, s, i0.String(), i1.String())
	Tei(t, s, i0.Int(), i1.Int())
}

func Tes(t *testing.T, s string, s0, s1 string) {
	if s0 != s1 {
		t.Fatalf("%s -> %s ≠ %s\n", s, s0, s1)
	}
}

func Tei(t *testing.T, s string, s0, s1 int64) {
	if s0 != s1 {
		t.Fatalf("%s -> %d ≠ %d\n", s, s0, s1)
	}
}
