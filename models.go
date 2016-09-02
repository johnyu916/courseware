package main

import (
	"github.com/attic-labs/noms/go/dataset"
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

func getSubmissions(ds dataset.Dataset) types.Map {
	hv, ok := ds.MaybeHeadValue()
	if ok {
		return hv.(types.Map)
	} else {
		return types.NewMap()
	}
}
