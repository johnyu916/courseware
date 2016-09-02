package courseware

import (
	"fmt"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/dataset"
	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
)

type Exam struct {
	Name     string
	Problems []Problem
}

type Problem struct {
	Id     uint64
	Text   string
	Answer string
}

type Submission struct {
	UserId    uint64
	ProblemId uint64
	Submitted string
}

func GetSubmissions(ds dataset.Dataset) types.Map {
	hv, ok := ds.MaybeHeadValue()
	if ok {
		return hv.(types.Map)
	} else {
		return types.NewMap()
	}
}

func getHistory() []Submission {
	dbpath := "http://localhost:8000::courseware"
	database, value, err := spec.GetPath(dbpath)
	defer database.Close()
	if err != nil {
		fmt.Println("getting database from noms failed")
		panic(err)
	}

	origCommit, _ := value.(types.Struct)
	iter := NewCommitIterator(database, origCommit)
	var submissions []Submission
	for node, ok := iter.Next(); ok; node, ok = iter.Next() {
		hashStr := node.commit.Hash().String()
		fmt.Println(hashStr)

		value := node.commit.Get(datas.ValueField)
		valMap := value.(types.Map)
		valMap.IterAll(func(k, v types.Value) {
			var s Submission
			err := marshal.Unmarshal(v, &s)
			if err != nil {
				fmt.Println("Error during unmarshalling")
			} else {
				fmt.Println(s)
			}

		})
	}
	return submissions
}
