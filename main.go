package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	. "srv2/utils"
)

const MaxIters = 30000

var TemplatePath = path.Join("templates", "index.html")

var tasksSpawner = []*Task{
	NewTask(20000, 3000),
	NewTask(5000, 1000),
	NewTask(16000, 2000),
	NewTask(10000, 2000),
	NewTask(15000, 500),
	NewTask(12000, 1000),
}

var RM = func(pq *PriorityQueue) func(i int, j int) bool {
	return func(i int, j int) bool {
		t1 := pq.Get(i).(*Task)
		t2 := pq.Get(j).(*Task)

		if t1.Id() == t2.Id() {
			return t1.ExecTimeRemaining() < t2.ExecTimeRemaining()
		}
		return t1.Period() < t2.Period()

		//return t1.Count * t1.Period() < t2.Count * t2.Period()
	}
}

func summaryLoad(tasks []*Task) float64 {
	sum := .0
	for _, task := range tasks {
		fmt.Println(float64(task.ExecTime()) / float64(task.Period()))
		sum += float64(task.ExecTime()) / float64(task.Period())
	}
	return sum
}

func handler(w http.ResponseWriter, _ *http.Request) {

	pq := NewPriorityQueue(RM)

	traceData := make([]map[string]interface{}, 0)
	tasks := make([]Task, 0)

	for moment := uint64(0); moment < MaxIters; moment++ {
		for _, task := range tasksSpawner {
			inst := NewTask(task.Period(), task.ExecTime())
			inst.SetName(task.Name())
			inst.Spawn(moment, pq, func() {
				task.Count++
				inst.Count = task.Count
			})
		}

		if t := pq.Peek(); t != nil {
			t.(*Task).Execute(moment, pq, func(task *Task) {
				tasks = append(tasks, *task)
			})
		}
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

	for _, task := range tasks {
		periods := make([]map[string]interface{}, 0)
		markers := make([]map[string]interface{}, 0)

		for _, marker := range task.Markers() {
			markerData := map[string]interface{}{
				"type":  "diamond",
				"value": marker.Time,
			}
			if !marker.IsStart {
				markerData["fill"] = "#FF0000"
			}
			markers = append(markers, markerData)
		}

		for _, execMoment := range task.ExecMoments() {
			periods = append(periods, map[string]interface{}{
				"start": execMoment,
				"end":   execMoment,
			})
		}

		traceData = append(traceData, map[string]interface{}{
			"id":      task.Id(),
			"name":    task.Name(),
			"markers": markers,
			"periods": periods,
		})
	}

	if err := tmpl.Execute(w, traceData); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("Summary load: %f\n", summaryLoad(tasksSpawner))
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	log.Print("Start server http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
