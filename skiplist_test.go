package skiplist

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/itsmontoya/mappedslice"
	"github.com/itsmontoya/rbt"
)

var (
	testSet_100_000     = getTestSet(100_000)
	testSet_100_000_000 = getTestSet(100_000_000)
)

func TestSkiplist_Set(t *testing.T) {
	type testcase struct {
		name string
		es   []Entry[Varchar16, int]
	}

	testcases := []testcase{
		{
			name: "basic",
			es: []Entry[Varchar16, int]{
				{
					Key:   MakeVarchar16("0001"),
					Value: 1,
				},
				{
					Key:   MakeVarchar16("0002"),
					Value: 2,
				},
				{
					Key:   MakeVarchar16("0003"),
					Value: 3,
				},
				{
					Key:   MakeVarchar16("0004"),
					Value: 4,
				},
				{
					Key:   MakeVarchar16("0005"),
					Value: 5,
				},
				{
					Key:   MakeVarchar16("0006"),
					Value: 6,
				},
				{
					Key:   MakeVarchar16("0007"),
					Value: 7,
				},
				{
					Key:   MakeVarchar16("0008"),
					Value: 8,
				},
			},
		},
		{
			name: "reverse",
			es: []Entry[Varchar16, int]{
				{
					Key:   MakeVarchar16("0008"),
					Value: 8,
				},
				{
					Key:   MakeVarchar16("0007"),
					Value: 7,
				},
				{
					Key:   MakeVarchar16("0006"),
					Value: 6,
				},
				{
					Key:   MakeVarchar16("0005"),
					Value: 5,
				},
				{
					Key:   MakeVarchar16("0004"),
					Value: 4,
				},
				{
					Key:   MakeVarchar16("0003"),
					Value: 3,
				},
				{
					Key:   MakeVarchar16("0002"),
					Value: 2,
				},
				{
					Key:   MakeVarchar16("0001"),
					Value: 1,
				},
			},
		},
		{
			name: "split",
			es: []Entry[Varchar16, int]{
				{
					Key:   MakeVarchar16("0005"),
					Value: 5,
				},
				{
					Key:   MakeVarchar16("0006"),
					Value: 6,
				},
				{
					Key:   MakeVarchar16("0007"),
					Value: 7,
				},
				{
					Key:   MakeVarchar16("0008"),
					Value: 8,
				},
				{
					Key:   MakeVarchar16("0001"),
					Value: 1,
				},
				{
					Key:   MakeVarchar16("0002"),
					Value: 2,
				},
				{
					Key:   MakeVarchar16("0003"),
					Value: 3,
				},
				{
					Key:   MakeVarchar16("0004"),
					Value: 4,
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer os.RemoveAll("testing")
			sl, err := New[Varchar16, int](tc.name, "testing", 32)
			if err != nil {
				t.Fatal(err)
			}

			sl.incrementEvery = 2
			for _, e := range tc.es {
				if err := sl.Set(&e.Key, &e.Value); err != nil {
					t.Fatal(err)
				}

				sl.printTree()
			}

			//sl.printTree()

			val, _ := sl.Get(MakeVarchar16("0008"))
			fmt.Println("0008", val)
			val, err = sl.Get(MakeVarchar16("0007"))
			fmt.Println("0007", val, err)
		})
	}
}

func BenchmarkSkiplist_Set_100_000(b *testing.B) {
	benchmarkSkiplist_Set(b, testSet_100_000)
}

func BenchmarkSkiplist_Set_100_000_000(b *testing.B) {
	benchmarkSkiplist_Set(b, testSet_100_000_000)
}

func BenchmarkBoltDB_Put_100_000(b *testing.B) {
	benchmarkBoltDB_Put(b, testSet_100_000)
}

func BenchmarkBoltDB_Put_100_000_000(b *testing.B) {
	benchmarkBoltDB_Put(b, testSet_100_000_000)
}

func BenchmarkRBT_Put_100_000(b *testing.B) {
	benchmarkRBT_Put(b, testSet_100_000)
}

func BenchmarkRBT_Put_100_000_000(b *testing.B) {
	benchmarkRBT_Put(b, testSet_100_000_000)
}

func BenchmarkMappedSlice_Append(b *testing.B) {
	defer os.RemoveAll("testing")
	if err := os.MkdirAll("./testing", 0744); err != nil {
		b.Fatal(err)
	}

	s, err := mappedslice.New[Entry[Varchar32, int]]("./testing/mapped.bat", 1024*1024)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer s.Close()
	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		if j >= len(testSet_100_000_000) {
			j = 0
		}

		e := testSet_100_000_000[j]
		if err = s.Append(e); err != nil {
			b.Fatal(err)
			return
		}

		j++
	}

	b.ReportAllocs()
}

func benchmarkSkiplist_Set(b *testing.B, testSet []Entry[Varchar32, int]) {
	defer os.RemoveAll("testing")
	s, err := New[Varchar32, int]("insert_benchmark", "testing", 1024*1024)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer s.Close()
	b.ResetTimer()

	var j int
	for i := 0; i < b.N; i++ {
		if j >= len(testSet) {
			j = 0
		}

		e := testSet[j]
		if err = s.Set(&e.Key, &e.Value); err != nil {
			b.Fatal(err)
			return
		}

		j++
	}

	b.ReportAllocs()
}

func benchmarkBoltDB_Put(b *testing.B, testSet []Entry[Varchar32, int]) {
	defer os.RemoveAll("testing")

	db, err := bolt.Open("testing", 0744, nil)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer db.Close()

	var j int
	if err = db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte("testing"))
		if err != nil {
			return err
		}

		for i := 0; i < b.N; i++ {
			if j >= len(testSet) {
				j = 0
			}

			e := testSet[j]
			if err = bkt.Put(e.Key[:], []byte(strconv.Itoa(e.Value))); err != nil {
				return err
			}

			j++
		}

		return nil
	}); err != nil {
		b.Fatal(err)
		return
	}

	b.ReportAllocs()
}

func benchmarkRBT_Put(b *testing.B, testSet []Entry[Varchar32, int]) {
	defer os.RemoveAll("testing")
	if err := os.MkdirAll("testing", 0744); err != nil {
		b.Fatal(err)
		return
	}

	db, err := rbt.NewMMAP("testing", "rbt", 1_000_000)
	if err != nil {
		b.Fatal(err)
		return
	}
	defer db.Close()

	var j int
	for i := 0; i < b.N; i++ {
		if j >= len(testSet) {
			j = 0
		}

		e := testSet[j]
		db.Put(e.Key[:], []byte(strconv.Itoa(e.Value)))
		j++
	}

	b.ReportAllocs()
}

func getTestSet(n int) (set []Entry[Varchar32, int]) {
	set = make([]Entry[Varchar32, int], 0, n)
	for i := 0; i < n; i++ {
		str := fmt.Sprintf("'%032dkm'", i)
		key := MakeVarchar32(str)
		e := makeEntry(&key, &i)
		set = append(set, e)
	}

	return
}
