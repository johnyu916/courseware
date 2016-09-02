package main

import (
	"encoding/json"
	"fmt"
	"github.com/attic-labs/noms/go/dataset"
	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"html/template"
	"io/ioutil"
	"net/http"
)

var exam Exam
var ds dataset.Dataset

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// load a problem
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func problemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("a GET request to problem")
	} else {
		fmt.Println("a POST request to problem")
		value := r.FormValue("answer")
		// add this to DB
		submission := Submission{0, 0, value}
		np, err := marshal.Marshal(submission)
		if err != nil {
			fmt.Println("There was a problem trying to marshal the submisison")
			fmt.Println(err)
			return
		}
		submissions := getSubmissions(ds)
		submissions = submissions.Set(types.Number(0), np)
		_, err = ds.CommitValue(submissions)
		if err != nil {
			fmt.Println("There was a problem trying to commit to db")
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully committed submission")
	}
	t, _ := template.ParseFiles("problem.html")
	t.Execute(w, exam.Problems[0])

}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	// show the most recent version
	var submissions []Submission
	data := getSubmissions(ds)
	data.IterAll(func(k, v types.Value) {
		var s Submission
		err := marshal.Unmarshal(v, &s)
		if err == nil {
			submissions = append(submissions, s)
		} else {
			fmt.Println("Error unmarshalling submission")
			fmt.Println(err)
		}
	})
	t, _ := template.ParseFiles("history.html")
	t.Execute(w, submissions[0])
}

func main() {
	// read the exam
	text, err := ioutil.ReadFile("exam.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal([]byte(text), &exam)

	ds, err = spec.GetDataset("http://localhost:8000::courseware")
	if err != nil {
		fmt.Println("getting dataset from noms failed")
		panic(err)
	}

	fmt.Println(exam)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/problem", problemHandler)
	http.HandleFunc("/history", historyHandler)
	http.ListenAndServe(":8080", nil)
}
