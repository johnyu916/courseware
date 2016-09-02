package main

import (
	"bytes"
	"fmt"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"github.com/johnyu916/courseware"
	"io"
)

func main() {
	database, value, err := spec.GetPath("http://localhost:8000::courseware")
	fmt.Println(database)
	fmt.Println(value)
	fmt.Println(err)
	defer database.Close()

	origCommit, ok := value.(types.Struct)
	if !ok {
		fmt.Println("Something is wrong")
		return
	}

	iter := courseware.NewCommitIterator(database, origCommit)
	for ln, ok := iter.Next(); ok; ln, ok = iter.Next() {
		buff := &bytes.Buffer{}
		printCommit(ln, buff, database)
	}
}

func printCommit(node LogNode, w io.Writer, db datas.Database) {
	hashStr := node.commit.Hash().String()
	fmt.Println(hashStr)
	//parents := commitRefsFromSet(node.commit.Get(datas.ParentsField).(types.Set))
	value := node.commit.Get(datas.ValueField)
	valMap := value.(types.Map)
	valMap.IterAll(func(k, v types.Value) {
		var s courseware.Submission
		err := marshal.Unmarshal(v, &s)
		if err != nil {
			fmt.Println("Error during unmarshalling")
		} else {
			fmt.Println(s)
		}
	})
}
