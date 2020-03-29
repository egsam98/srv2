package main

import (
	"encoding/json"
	"fmt"
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
	NewTask(1, 60000, 2500),
	NewTask(2, 10000, 1500),
	NewTask(3, 12000, 3000),
	NewTask(4, 12000, 2000),
	NewTask(5, 6000, 500),
	NewTask(6, 20000, 1000),
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	log.Print(fmt.Sprintf("Start server http://localhost:%s...", PORT))
	err := http.ListenAndServe(":"+PORT, nil)
	Panic(err)
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
	title, traceData, err := services.NewSchedullingService(TasksConfig).Run(method)
	Panic(err)

	tmpl, err := template.ParseFiles(TemplatePath)
	Panic(err)

	bytes, err := json.Marshal(traceData)
	Panic(err)
	err = ioutil.WriteFile(OUT, bytes, os.ModePerm)
	Panic(err)

	responseBody := map[string]interface{}{
		"title":     title,
		"traceData": traceData,
	}
	err = tmpl.Execute(w, responseBody)
	Panic(err)
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}
