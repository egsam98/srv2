package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"srv2/utils"
)

const MaxIters = 20000

var TemplatePath = path.Join("templates", "index.html")

var RM = func(pq *utils.PriorityQueue) func(i int, j int) bool {
	return func(i int, j int) bool {
		t1 := pq.Get(i).(*utils.Task)
		t2 := pq.Get(j).(*utils.Task)
		return t1.Period() < t2.Period()
	}
}

type TraceData []map[string]interface{}

func (td *TraceData) find(task *utils.Task) (map[string]interface{}, bool) {
	for _, elem := range *td {
		if elem["id"] == task.Id() {
			return elem, true
		}
	}
	return nil, false
}

func (td *TraceData) Add(moment uint64, task *utils.Task) {
	fmt.Printf("moment: %d, %+v\n", moment, *task)

	if found, exists := td.find(task); exists {
		found["periods"] = task.Timelines()
	} else {
		*td = append(*td, map[string]interface{}{
			"id":      task.Id(),
			"name":    task.Name(),
			"periods": task.Timelines(),
		})
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	var tasks = []*utils.Task{
		//utils.NewTask(20000, 3000),
		//utils.NewTask(5000, 2000),
		//utils.NewTask(10000, 2000),
		utils.NewTask(2000, 1000),
		utils.NewTask(3000, 1500),
	}

	traceData := make(TraceData, 0)

	pq := utils.NewPriorityQueue(RM)

	var moment uint64 = 0
	for moment < MaxIters {
		for _, task := range tasks {
			task.Spawn(moment, pq, traceData.Add)
		}

		if task := pq.Peek(); task != nil {
			task.(*utils.Task).Execute(moment, pq, traceData.Add)
		}

		moment++
	}

	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), 500)
			log.Fatal(err)
		}
	}()

	tmpl, err := template.ParseFiles(TemplatePath)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(w, traceData); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	log.Print("Start server http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
