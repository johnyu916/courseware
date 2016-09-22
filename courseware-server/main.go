package main

import (
	"encoding/json"
	"fmt"
	"github.com/attic-labs/noms/go/dataset"
	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
	"github.com/johnyu916/courseware"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Exam courseware.Exam
type Submission courseware.Submission

var exam Exam
var ds dataset.Dataset

type History struct {
	//Submissions []Submission
	Commits []courseware.Commit
}

type ProblemContext struct {
	ProblemSubmissions []ProblemSubmission
}

type ProblemSubmission struct {
	Problem    courseware.Problem
	Submission Submission
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// load a problem
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}

func problemHandler(w http.ResponseWriter, r *http.Request) {
	submissionsMap := courseware.GetSubmissions(ds)
	if r.Method == "GET" {
		fmt.Println("a GET request to problem")
	} else {
		fmt.Println("a POST request to problem")
		value := r.FormValue("answer")
		ids := r.FormValue("problemid")
		// add this to DB
		id, err := strconv.ParseUint(ids, 10, 64)
		if err != nil {
			panic("Cannot convert to uint" + err.Error())
		}
		submission := Submission{0, id, value}
		np, err := marshal.Marshal(submission)
		if err != nil {
			fmt.Println("There was a problem trying to marshal the submisison")
			fmt.Println(err)
			return
		}
		submissionsMap = submissionsMap.Set(types.Number(id), np)
		_, err = ds.CommitValue(submissionsMap)
		if err != nil {
			fmt.Println("There was a problem trying to commit to db")
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully committed submission")
	}
	submissionsMap = courseware.GetSubmissions(ds)
	var submissions []Submission
	submissionsMap.IterAll(func(k, v types.Value) {
		var s Submission
		err := marshal.Unmarshal(v, &s)
		if err != nil {
			panic("Error during problemHandler unmarshalling")
		} else {
			submissions = append(submissions, s)
		}
	})
	var pss []ProblemSubmission
	for ind, problem := range exam.Problems {
		ps := ProblemSubmission{problem, submissions[ind]}
		pss = append(pss, ps)
	}
	t, _ := template.ParseFiles("problem.html")
	t.Execute(w, ProblemContext{pss})

}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	// show the most recent version
	commits := courseware.GetHistory()
	/*
		var submissions []Submission
			data := courseware.GetSubmissions(ds)
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
	*/
	t, _ := template.ParseFiles("history.html")
	//t.Execute(w, History{submissions})
	t.Execute(w, History{commits})
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
