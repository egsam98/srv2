package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"srv2/services"
	. "srv2/utils"
	"strings"
)

const PORT = "8080"
const OUT = "out.json"

var TemplatePath = path.Join("templates", "index.html")
var TasksConfig = []*Task{
	NewTask(1, 20000, 3000),
	NewTask(2, 5000, 1000),
	NewTask(3, 16000, 3000),
	NewTask(4, 13000, 2000),
	NewTask(5, 15000, 500),
	NewTask(6, 12000, 1000),
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	log.Print("Start server http://localhost:8080...")
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			errBody := map[string]string{"error": err.(error).Error()}
			errBodyByteArr, err := json.Marshal(errBody)
			if err != nil {
				log.Fatal(err)
			}
			http.Error(w, string(errBodyByteArr), 500)
		}
	}()

	method := strings.TrimPrefix(r.URL.Path, "/")
	title, traceData := services.NewSchedullingService(TasksConfig).Run(method)

	tmpl, err := template.ParseFiles(TemplatePath)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(traceData)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(OUT, bytes, os.ModePerm)
	if err != nil {
		panic(err)
	}

	responseBody := map[string]interface{}{
		"title":     title,
		"traceData": traceData,
	}
	if err := tmpl.Execute(w, responseBody); err != nil {
		panic(err)
	}
}
