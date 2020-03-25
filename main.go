package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"srv2/services"
	"strings"
)

var TemplatePath = path.Join("templates", "index.html")

func handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, err.(error).Error(), 500)
			//log.Fatal(err)
		}
	}()

	method := strings.TrimPrefix(r.URL.Path, "/")
	title, traceData := services.NewSchedullingService().Run(method)

	tmpl, err := template.ParseFiles(TemplatePath)
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

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	log.Print("Start server http://localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
