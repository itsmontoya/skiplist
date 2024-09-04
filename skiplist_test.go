package skiplist

import (
	"fmt"
	"testing"
)

func TestSkiplist_Insert(t *testing.T) {

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
			//defer os.RemoveAll(path.Join("testing", tc.name))
			sl, err := New[Varchar16, int](tc.name, "testing")
			if err != nil {
				t.Fatal(err)
			}

			sl.incrementEvery = 2
			for _, e := range tc.es {
				if err := sl.Set(e.Key, e.Value); err != nil {
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
